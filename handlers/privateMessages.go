package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return uuid.Nil, err
	}
	sessionToken := cookie.Value
	return db.GetUserIDFromSession(sessionToken)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverID, err := uuid.FromString(r.FormValue("receiver_id"))
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

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := uuid.FromString(r.FormValue("message_id"))
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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsersOrderedByLastMessageOrAlphabetically()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
