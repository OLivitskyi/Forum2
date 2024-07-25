package handlers

import (
	"forum/db"
	"log"
	"net/http"
	"strconv"
	"sync"

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
	}
	defer ws.Close()

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		return
	}
	clients[ws] = userID

	for {
		var msg db.WebSocketMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			if clients[client] == msg.Receiver {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Error writing to websocket: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func WebSocketHandler() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
}
