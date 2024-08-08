// handlers/websocket.go
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

type PostMessage struct {
	MessageID uuid.UUID `json:"message_id"`
	Post      db.Post   `json:"post"`
	Timestamp time.Time `json:"timestamp"`
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

		if postType, ok := message["type"].(string); ok && postType == "post" {
			var post db.Post
			postData, _ := json.Marshal(message["data"])
			if err := json.Unmarshal(postData, &post); err != nil {
				log.Printf("Error unmarshalling post data: %v", err)
				continue
			}
			log.Printf("Received post from user %s: %v", userID, post)
			postMessage := PostMessage{
				MessageID: uuid.Must(uuid.NewV4()),
				Post:      post,
				Timestamp: time.Now(),
			}
			broadcast <- postMessage
		} else {
			log.Printf("Unknown message format: %v", message)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		switch m := msg.(type) {
		case PostMessage:
			log.Printf("Broadcasting new post: %+v", m.Post)
			mutex.Lock()
			if _, exists := processedMessages[m.MessageID]; exists {
				mutex.Unlock()
				continue
			}
			processedMessages[m.MessageID] = struct{}{}
			mutex.Unlock()

			completePost, err := db.GetPostByID(m.Post.ID)
			if err != nil {
				log.Printf("Error fetching complete post data: %v", err)
				continue
			}
			for client := range clients {
				err := client.WriteJSON(completePost)
				if err != nil {
					log.Printf("Error writing to websocket: %v", err)
					client.Close()
					mutex.Lock()
					delete(clients, client)
					mutex.Unlock()
				}
			}
		case db.WebSocketMessage:
			log.Printf("Broadcasting message: %v", m)
			err := db.AddMessage(m.Sender, m.Receiver, m.Content)
			if err != nil {
				log.Printf("Error storing message in the database: %v", err)
				continue
			}
			for client := range clients {
				if clients[client] == m.Receiver {
					err := client.WriteJSON(m)
					if err != nil {
						log.Printf("Error writing to websocket: %v", err)
						client.Close()
						mutex.Lock()
						delete(clients, client)
						mutex.Unlock()
					}
				}
			}
		case db.ReactionMessage:
			log.Printf("Broadcasting reaction: %v", m)
			if m.PostID != uuid.Nil {
				err := db.AddPostReaction(m.UserID, m.PostID, m.ReactionType)
				if err != nil {
					log.Printf("Error storing post reaction in the database: %v", err)
					continue
				}
			} else if m.CommentID != uuid.Nil {
				err := db.AddCommentReaction(m.UserID, m.CommentID, m.ReactionType)
				if err != nil {
					log.Printf("Error storing comment reaction in the database: %v", err)
					continue
				}
			}
			for client := range clients {
				err := client.WriteJSON(m)
				if err != nil {
					log.Printf("Error writing to websocket: %v", err)
					client.Close()
					mutex.Lock()
					delete(clients, client)
					mutex.Unlock()
				}
			}
		default:
			log.Printf("Unknown message type: %T", m)
		}
	}
}

func InitWebSocketHandler() {
	go handleMessages()
}
