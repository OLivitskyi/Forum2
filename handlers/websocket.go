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
		var msg interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}
		switch m := msg.(type) {
		case map[string]interface{}:
			if postType, ok := m["type"].(string); ok && postType == "post" {
				var post db.Post
				postData, _ := json.Marshal(m["data"])
				if err := json.Unmarshal(postData, &post); err != nil {
					log.Printf("Error unmarshalling post data: %v", err)
					continue
				}
				log.Printf("Received post from user %s: %v", userID, post)
				broadcast <- post
			} else {
				log.Printf("Unknown map message type: %v", m)
			}
		case db.WebSocketMessage:
			m.Timestamp = time.Now().Format(time.RFC3339)
			log.Printf("Received message from user %s: %v", userID, m)
			broadcast <- m
		case db.ReactionMessage:
			m.UserID = userID
			log.Printf("Received reaction from user %s: %v", userID, m)
			broadcast <- m
		case db.Post:
			log.Printf("Received post from user %s: %v", userID, m)
			broadcast <- m
		default:
			log.Printf("Unknown message type: %T", m)
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		switch m := msg.(type) {
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
		case db.Post:
			log.Printf("Broadcasting new post: %+v", m)
			completePost, err := db.GetPostByID(m.ID)
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
		default:
			log.Printf("Unknown message type: %T", m)
		}
	}
}

func InitWebSocketHandler() {
	go handleMessages()
}
