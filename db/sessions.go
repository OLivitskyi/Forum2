package db

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// SaveSession saves a new session in the database.
func SaveSession(token string, userID uuid.UUID, expiration time.Time) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token, userID, expiration)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSession deletes a session from the database.
func DeleteSession(token string) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`DELETE FROM sessions WHERE token = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token)
	if err != nil {
		return err
	}
	return nil
}

// GetSessionByUserID retrieves a session by user ID.
func GetSessionByUserID(userID uuid.UUID) (*Session, error) {
	var session Session
	err := DB.QueryRow(`SELECT token, user_id, expires_at FROM sessions WHERE user_id = ?`, userID).Scan(
		&session.SessionToken, &session.UserID, &session.ExpireTime)
	if err != nil {
		log.Printf("Error finding session for user ID %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Session token found for user ID: %s", session.UserID)
	return &session, nil
}

// GetUserIDFromSession retrieves the user ID associated with a session token.
func GetUserIDFromSession(token string) (uuid.UUID, error) {
	if DB == nil {
		return uuid.Nil, fmt.Errorf("db connection failed")
	}
	var userID uuid.UUID
	err := DB.QueryRow(`SELECT user_id FROM sessions WHERE token = ?`, token).Scan(&userID)
	if err != nil {
		log.Printf("Error finding session token: %v", err)
		return uuid.Nil, err
	}
	log.Printf("Session token found for user ID: %s", userID)
	return userID, nil
}

// ClearSessions deletes all sessions from the database.
func ClearSessions() error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	_, err := DB.Exec("DELETE FROM sessions")
	if err != nil {
		return fmt.Errorf("failed to clear sessions: %v", err)
	}
	log.Println("All sessions have been cleared from the database")
	return nil
}
