package handlers

import (
	"forum/db"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan db.WebSocketMessage)
var mutex = &sync.Mutex{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error upgrading to websocket: %v", err)
		return
	}
	defer func() {
		ws.Close()
		log.Printf("User %d disconnected", clients[ws])
		mutex.Lock()
		delete(clients, ws)
		mutex.Unlock()
		db.UpdateUserStatus(clients[ws], false)
	}()

	// Read userID from session
	userID, err := getUserIDFromSession(r)
	if err != nil {
		log.Printf("Unauthorized access: %v", err)
		return
	}
	log.Printf("User %d connected", userID)
	mutex.Lock()
	clients[ws] = userID
	mutex.Unlock()
	db.UpdateUserStatus(userID, true)

	for {
		var msg db.WebSocketMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}
		msg.Timestamp = time.Now().Format(time.RFC3339)
		log.Printf("Received message from user %d: %v", userID, msg)
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		log.Printf("Broadcasting message: %v", msg)
		// Store message in the database
		err := db.AddMessage(msg.Sender, msg.Receiver, msg.Content)
		if err != nil {
			log.Printf("Error storing message in the database: %v", err)
			continue
		}
		// Broadcast message to all connected clients
		for client := range clients {
			if clients[client] == msg.Receiver {
				err := client.WriteJSON(msg)
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
}

func WebSocketHandler() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
}
