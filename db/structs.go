package db

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
}

type Login struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	MessageID  uuid.UUID `json:"message_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	IsRead     bool      `json:"is_read"`
}

type UserStatus struct {
	UserID       uuid.UUID `json:"user_id"`
	IsOnline     bool      `json:"is_online"`
	LastActivity time.Time `json:"last_activity"`
}

type WebSocketMessage struct {
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Sender    uuid.UUID `json:"sender"`
	Receiver  uuid.UUID `json:"receiver"`
	Timestamp string    `json:"timestamp"`
	IsRead    bool      `json:"is_read"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	User         *User       `json:"user,omitempty"`
	Subject      string      `json:"subject"`
	Content      string      `json:"content"`
	Categories   []*Category `json:"categories,omitempty"`
	Comments     []*Comment  `json:"comments,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	LikeCount    int         `json:"like_count"`
	DislikeCount int         `json:"dislike_count"`
}

type Comment struct {
	ID           uuid.UUID `json:"id"`
	PostID       uuid.UUID `json:"post_id"`
	UserID       uuid.UUID `json:"user_id"`
	User         *User     `json:"user,omitempty"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
}

type PostCategory struct {
	PostID     uuid.UUID `json:"post_id"`
	CategoryID int       `json:"category_id"`
}

type Reaction struct {
	UserID     uuid.UUID    `json:"user_id"`
	PostID     uuid.UUID    `json:"post_id,omitempty"`
	CommentID  uuid.UUID    `json:"comment_id,omitempty"`
	ReactionID int          `json:"reaction_id"`
	Type       ReactionType `json:"type"`
}

type ReactionType string

const (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
)

type Session struct {
	Username     string    `json:"username"`
	SessionToken string    `json:"session_token"`
	ExpireTime   time.Time `json:"expire_time"`
}
