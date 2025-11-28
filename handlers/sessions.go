package handlers

import (
	"forum/db"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

var sessions = map[string]db.Session{}

// NewSession creates a new session for the user
func NewSession(w http.ResponseWriter, username string, userID uuid.UUID) (string, error) {
	if isSessionUp(username) {
		return "", nil
	}
	token, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
		return "", err
	}
	session := db.Session{
		Username:     username,
		SessionToken: token.String(),
		ExpireTime:   time.Now().Add(100 * time.Minute),
	}
	sessions[token.String()] = session
	expiration := time.Now().Add(4 * time.Hour)
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	err = db.SaveSession(token.String(), userID, expiration)
	if err != nil {
		log.Fatalf("Failed to save session to database: %v", err)
		return "", err
	}
	log.Printf("New session created for user %s with token %s", username, token.String())
	return token.String(), nil
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
	token, err := getSessionToken(r)
	if err != nil {
		log.Println("Session token not found:", err)
		return ""
	}
	key, ok := sessions[token]
	if ok {
		return key.Username
	} else {
		log.Printf("Invalid session token: %s", token)
		return ""
	}
}

// SessionExpired checks if the session has expired
func SessionExpired(r *http.Request) bool {
	token, err := getSessionToken(r)
	if err != nil {
		log.Println("Session token not found:", err)
		return true
	}
	key, ok := sessions[token]
	if ok {
		return key.ExpireTime.Before(time.Now())
	}
	return true
}

// CloseSession closes the session and deletes the cookie
func CloseSession(w http.ResponseWriter, r *http.Request) {
	token, err := getSessionToken(r)
	if err != nil {
		http.Error(w, "No session token found", http.StatusBadRequest)
		return
	}
	log.Printf("Attempting to close session: %s", token)
	_, ok := sessions[token]
	if ok {
		delete(sessions, token)
		cookie := http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		}
		http.SetCookie(w, &cookie)
		err = db.DeleteSession(token)
		if err != nil {
			log.Printf("Failed to delete session from database: %v", err)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

// RequireLogin is a middleware that checks for a valid session
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if SessionExpired(r) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}
		username := ValidateSession(r)
		if username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("No session cookie found:", err)
		return uuid.Nil, err
	}
	sessionToken := cookie.Value
	log.Println("Session token received:", sessionToken)
	return db.GetUserIDFromSession(sessionToken)
}

// getSessionToken extracts the session token from the request
func getSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
