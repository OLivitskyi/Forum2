package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	Username     string
	sessionToken string
	expireTime   time.Time
}

var sessions = map[string]Session{}

// NewSession creates a new session for the user
func NewSession(w http.ResponseWriter, username string) {
	if isSessionUp(username) {
		return
	}

	token, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}

	session := Session{
		Username:     username,
		sessionToken: token.String(),
		expireTime:   time.Now().Add(100 * time.Minute),
	}

	sessions[token.String()] = session

	expiration := time.Now().Add(4 * time.Hour)
	cookie := http.Cookie{
		Name:    "session",
		Value:   session.sessionToken,
		Expires: expiration,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
}

// isSessionUp checks if a session is active for the user
func isSessionUp(username string) bool {
	for _, a := range sessions {
		if a.Username == username {
			return true
		}
	}
	return false
}

// ValidateSession validates the session from the request
func ValidateSession(r *http.Request) string {
	sessionToken, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	key, ok := sessions[sessionToken.Value]
	if ok {
		return key.Username
	} else {
		return ""
	}
}

// SessionExpired checks if the session has expired
func SessionExpired(r *http.Request) bool {
	sessionToken, err := r.Cookie("session")
	if err != nil {
		return true
	}
	key, ok := sessions[sessionToken.Value]
	if ok {
		return key.expireTime.Before(time.Now())
	}
	return true
}

// CloseSession closes the session and deletes the cookie
func CloseSession(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := r.Cookie("session")
	if err != nil {
		return
	}

	_, ok := sessions[sessionToken.Value]
	if ok {
		delete(sessions, sessionToken.Value)
		cookie := http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		}
		http.SetCookie(w, &cookie)
	}
}
