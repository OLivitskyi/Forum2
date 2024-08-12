package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

// SendMessageHandler handles sending a message.
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData struct {
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	receiverID, err := uuid.FromString(requestData.ReceiverID)
	if err != nil {
		log.Println("Invalid receiver ID:", err)
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	err = db.AddMessage(senderID, receiverID, requestData.Content)
	if err != nil {
		log.Println("Failed to send message:", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetMessagesHandler handles fetching messages.
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	otherUserID, err := uuid.FromString(r.URL.Query().Get("user_id"))
	if err != nil {
		log.Println("Invalid user ID:", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
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

	messages, err := db.GetMessages(userID, otherUserID, limit, offset)
	if err != nil {
		log.Println("Failed to get messages:", err)
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// UpdateStatusHandler handles updating user status.
func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData struct {
		IsOnline bool `json:"is_online"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err = db.UpdateUserStatus(userID, requestData.IsOnline)
	if err != nil {
		log.Println("Failed to update status:", err)
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserStatusHandler handles fetching user statuses.
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

// MarkMessageAsReadHandler handles marking a message as read.
func MarkMessageAsReadHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData struct {
		MessageID string `json:"message_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	messageID, err := uuid.FromString(requestData.MessageID)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = db.MarkMessageAsRead(messageID, userID)
	if err != nil {
		http.Error(w, "Failed to mark message as read", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUsersHandler handles fetching users.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	users, err := db.GetUsersOrderedByLastMessageOrAlphabetically(userID)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
