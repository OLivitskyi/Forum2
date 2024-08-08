package db

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// AddMessage adds a new message to the database.
func AddMessage(senderID, receiverID uuid.UUID, content string) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO messages (message_id, sender_id, receiver_id, content, created_at, is_read) VALUES (?, ?, ?, ?, ?, 0)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	messageID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(messageID, senderID, receiverID, content, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// GetMessages retrieves messages between two users with pagination support.
func GetMessages(senderID, receiverID uuid.UUID, limit, offset int) ([]Message, error) {
	if DB == nil {
		return nil, fmt.Errorf("db connection failed")
	}
	rows, err := DB.Query(`SELECT message_id, sender_id, receiver_id, content, created_at, is_read FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?`, senderID, receiverID, receiverID, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.IsRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// MarkMessageAsRead marks a message as read by updating its status.
func MarkMessageAsRead(messageID uuid.UUID, userID uuid.UUID) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`UPDATE messages SET is_read = 1 WHERE message_id = ? AND receiver_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(messageID, userID)
	if err != nil {
		return err
	}
	return nil
}
