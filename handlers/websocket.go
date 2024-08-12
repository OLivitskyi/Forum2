package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]uuid.UUID)
var broadcast = make(chan interface{})
var mutex = &sync.Mutex{}
var processedMessages = make(map[uuid.UUID]struct{})

type Message interface {
	Process()
}

type PostMessage struct {
	MessageID uuid.UUID `json:"message_id"`
	Type      string    `json:"type"`
	Post      db.Post   `json:"post"`
	Timestamp time.Time `json:"timestamp"`
}

func (pm PostMessage) Process() {
	log.Printf("Broadcasting new post: %+v", pm.Post)
	if isProcessed(pm.MessageID) {
		return
	}
	completePost, err := db.GetPostByID(pm.Post.ID)
	if err != nil {
		log.Printf("Error fetching complete post data: %v", err)
		return
	}
	broadcastMessageToAll("post", completePost)
}

type CommentMessage struct {
	MessageID uuid.UUID  `json:"message_id"`
	Type      string     `json:"type"`
	Comment   db.Comment `json:"comment"`
	Timestamp time.Time  `json:"timestamp"`
}

func (cm CommentMessage) Process() {
	log.Printf("Broadcasting new comment: %+v", cm.Comment)
	if isProcessed(cm.MessageID) {
		return
	}
	completeComment, err := db.GetCommentByID(cm.Comment.ID)
	if err != nil {
		log.Printf("Error fetching complete comment data: %v", err)
		return
	}
	broadcastMessageToAll("comment", completeComment)
}

type PrivateMessage struct {
	MessageID uuid.UUID `json:"message_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func (pm PrivateMessage) Process() {
	if isProcessed(pm.MessageID) {
		return
	}
	broadcastPrivateMessageToClient(pm.ReceiverID, pm)
}


func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer func() {
		userID := clients[ws]
		mutex.Lock()
		delete(clients, ws)
		mutex.Unlock()
		db.UpdateUserStatus(userID, false)
		ws.Close()
		log.Printf("User %s disconnected", userID)
		broadcastUserStatus() // Broadcast user status when someone disconnects
	}()

	sessionToken := r.URL.Query().Get("session_token")
	if sessionToken == "" {
		log.Println("Unauthorized access: session token missing")
		return
	}
	userID, err := db.GetUserIDFromSession(sessionToken)
	if err != nil {
		log.Printf("Unauthorized access: %v", err)
		return
	}
	log.Printf("User %s connected", userID)
	mutex.Lock()
	clients[ws] = userID
	mutex.Unlock()
	db.UpdateUserStatus(userID, true)
	broadcastUserStatus() // Broadcast user status when someone connects

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		log.Printf("Received message: %v", message)

		if messageType, ok := message["type"].(string); ok {
			switch messageType {
			case "post":
				var post db.Post
				postData, _ := json.Marshal(message["data"])
				if err := json.Unmarshal(postData, &post); err != nil {
					log.Printf("Error unmarshalling post data: %v", err)
					continue
				}
				log.Printf("Received post from user %s: %v", userID, post)
				postMessage := PostMessage{
					MessageID: uuid.Must(uuid.NewV4()),
					Type:      "post",
					Post:      post,
					Timestamp: time.Now(),
				}
				broadcast <- postMessage
			case "comment":
				var comment db.Comment
				commentData, _ := json.Marshal(message["data"])
				if err := json.Unmarshal(commentData, &comment); err != nil {
					log.Printf("Error unmarshalling comment data: %v", err)
					continue
				}

				if comment.ID == uuid.Nil {
					log.Printf("Invalid comment ID: %v", comment.ID)
					continue
				}

				log.Printf("Received comment from user %s: %v", userID, comment)
				commentMessage := CommentMessage{
					MessageID: uuid.Must(uuid.NewV4()),
					Type:      "comment",
					Comment:   comment,
					Timestamp: time.Now(),
				}
				broadcast <- commentMessage
			case "request_user_status":
				log.Printf("Received request for user status from user %s", userID)
				broadcastUserStatusToClient(ws)
			case "private_message":
				var privateMsg PrivateMessage
				messageData, _ := json.Marshal(message["data"])
				if err := json.Unmarshal(messageData, &privateMsg); err != nil {
					log.Printf("Error unmarshalling private message data: %v", err)
					continue
				}
				log.Printf("Received private message from user %s to user %s: %v", userID, privateMsg.ReceiverID, privateMsg.Content)
				privateMsg.MessageID = uuid.Must(uuid.NewV4())
				privateMsg.SenderID = userID
				privateMsg.Timestamp = time.Now()

				// Зберегти повідомлення в базу даних
				err = db.AddMessage(privateMsg.SenderID, privateMsg.ReceiverID, privateMsg.Content)
				if err != nil {
					log.Printf("Failed to save private message: %v", err)
					continue
				}

				broadcast <- privateMsg

			default:
				log.Printf("Unknown message type: %v", messageType)
			}
		} else {
			log.Printf("Message type not found in received message: %v", message)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		if message, ok := msg.(Message); ok {
			message.Process()
		} else {
			log.Printf("Unknown message type: %T", msg)
		}
	}
}

func InitWebSocketHandler() {
	go handleMessages()
}

type UserStatusMessage struct {
	MessageID uuid.UUID       `json:"message_id"`
	Type      string          `json:"type"`
	Users     []db.UserStatus `json:"users"`
	Timestamp time.Time       `json:"timestamp"`
}

func (usm UserStatusMessage) Process() {
	log.Printf("Broadcasting user status update")
	if isProcessed(usm.MessageID) {
		return
	}
	broadcastMessageToAll("user_status", usm.Users)
}

func broadcastUserStatus() {
	users, err := db.GetAllUsersWithStatus()
	if err != nil {
		log.Printf("Error fetching user status: %v", err)
		return
	}
	broadcastMessageToAll("user_status", users)
}

func broadcastUserStatusToClient(client *websocket.Conn) {
	users, err := db.GetAllUsersWithStatus()
	if err != nil {
		log.Printf("Error fetching user status: %v", err)
		return
	}
	broadcastMessageToClient(client, "user_status", users)
}

// New helper functions
func isProcessed(messageID uuid.UUID) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := processedMessages[messageID]; exists {
		return true
	}
	processedMessages[messageID] = struct{}{}
	return false
}

func broadcastMessageToAll(messageType string, data interface{}) {
	message := map[string]interface{}{
		"type": messageType,
		"data": data,
	}
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("Error writing to websocket: %v", err)
			client.Close()
			mutex.Lock()
			delete(clients, client)
			mutex.Unlock()
		}
	}
}

func broadcastMessageToClient(client *websocket.Conn, messageType string, data interface{}) {
	message := map[string]interface{}{
		"type": messageType,
		"data": data,
	}
	err := client.WriteJSON(message)
	if err != nil {
		log.Printf("Error writing to websocket: %v", err)
		client.Close()
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
	}
}

func broadcastPrivateMessageToClient(receiverID uuid.UUID, message PrivateMessage) {
	messageJSON := map[string]interface{}{
		"type": "private_message",
		"data": message,
	}
	for client, id := range clients {
		if id == receiverID {
			err := client.WriteJSON(messageJSON)
			if err != nil {
				log.Printf("Error writing to websocket: %v", err)
				client.Close()
				mutex.Lock()
				delete(clients, client)
				mutex.Unlock()
			}
		}
	}
}