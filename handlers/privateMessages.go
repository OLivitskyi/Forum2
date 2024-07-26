package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"strconv"
)

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID, err := strconv.Atoi(r.FormValue("sender_id"))
	if err != nil {
		log.Println("Invalid sender ID:", err)
		http.Error(w, "Invalid sender ID", http.StatusBadRequest)
		return
	}
	receiverID, err := strconv.Atoi(r.FormValue("receiver_id"))
	if err != nil {
		log.Println("Invalid receiver ID:", err)
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}
	content := r.FormValue("content")

	err = db.AddMessage(senderID, receiverID, content)
	if err != nil {
		log.Println("Failed to send message:", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	senderID, err := strconv.Atoi(r.URL.Query().Get("sender_id"))
	if err != nil {
		log.Println("Invalid sender ID:", err)
		http.Error(w, "Invalid sender ID", http.StatusBadRequest)
		return
	}
	receiverID, err := strconv.Atoi(r.URL.Query().Get("receiver_id"))
	if err != nil {
		log.Println("Invalid receiver ID:", err)
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	messages, err := db.GetMessages(senderID, receiverID, limit, offset)
	if err != nil {
		log.Println("Failed to get messages:", err)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Println("Invalid user ID:", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	isOnline := r.FormValue("is_online") == "true"

	err = db.UpdateUserStatus(userID, isOnline)
	if err != nil {
		log.Println("Failed to update status:", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	statuses, err := db.GetUserStatus()
	if err != nil {
		log.Println("Failed to get user statuses:", err)
		http.Error(w, "Failed to get user statuses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func MarkMessageAsReadHandler(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.Atoi(r.FormValue("message_id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = db.MarkMessageAsRead(messageID)
	if err != nil {
		http.Error(w, "Failed to mark message as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsersOrderedByLastMessageOrAlphabetically()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
