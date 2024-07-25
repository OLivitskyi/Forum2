package db

import "time"

type User struct {
	Id        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type Login struct {
	Username string
	Email    string
	Password string
}

type Message struct {
	MessageID  int       `json:"message_id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `json:"is_read"`
}

type UserStatus struct {
	UserID       int       `json:"user_id"`
	IsOnline     bool      `json:"is_online"`
	LastActivity time.Time `json:"last_activity"`
}

type WebSocketMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	Sender   int    `json:"sender"`
	Receiver int    `json:"receiver"`
}
