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
	mutex.Lock()
	if _, exists := processedMessages[pm.MessageID]; exists {
		mutex.Unlock()
		return
	}
	processedMessages[pm.MessageID] = struct{}{}
	mutex.Unlock()

	completePost, err := db.GetPostByID(pm.Post.ID)
	if err != nil {
		log.Printf("Error fetching complete post data: %v", err)
		return
	}
	for client := range clients {
		log.Printf("Sending post to client: %v", clients[client])
		message := map[string]interface{}{
			"type": "post",
			"data": completePost,
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
}

type CommentMessage struct {
	MessageID uuid.UUID  `json:"message_id"`
	Type      string     `json:"type"`
	Comment   db.Comment `json:"comment"`
	Timestamp time.Time  `json:"timestamp"`
}

func (cm CommentMessage) Process() {
	log.Printf("Broadcasting new comment: %+v", cm.Comment)
	mutex.Lock()
	if _, exists := processedMessages[cm.MessageID]; exists {
		mutex.Unlock()
		return
	}
	processedMessages[cm.MessageID] = struct{}{}
	mutex.Unlock()

	// Отримання коментаря за його ID
	completeComment, err := db.GetCommentByID(cm.Comment.ID)
	if err != nil {
		log.Printf("Error fetching complete comment data: %v", err)
		return
	}
	for client := range clients {
		log.Printf("Sending comment to client: %v", clients[client])
		message := map[string]interface{}{
			"type": "comment",
			"data": completeComment,
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

				// Важливо: переконайтеся, що comment.ID встановлено правильно
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
