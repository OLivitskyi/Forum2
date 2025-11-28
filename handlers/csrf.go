package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

// CSRFToken represents a CSRF token with expiration
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFStore stores CSRF tokens per session
type CSRFStore struct {
	tokens map[string]CSRFToken
	mu     sync.RWMutex
}

// Global CSRF store
var csrfStore = &CSRFStore{
	tokens: make(map[string]CSRFToken),
}

const csrfTokenLength = 32
const csrfTokenExpiry = 1 * time.Hour

// GenerateCSRFToken generates a new CSRF token for a session
func GenerateCSRFToken(sessionToken string) (string, error) {
	bytes := make([]byte, csrfTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(bytes)

	csrfStore.mu.Lock()
	csrfStore.tokens[sessionToken] = CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(csrfTokenExpiry),
	}
	csrfStore.mu.Unlock()

	return token, nil
}

// ValidateCSRFToken validates a CSRF token for a session
func ValidateCSRFToken(sessionToken, csrfToken string) bool {
	csrfStore.mu.RLock()
	stored, exists := csrfStore.tokens[sessionToken]
	csrfStore.mu.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(stored.ExpiresAt) {
		// Token expired, remove it
		csrfStore.mu.Lock()
		delete(csrfStore.tokens, sessionToken)
		csrfStore.mu.Unlock()
		return false
	}

	return stored.Token == csrfToken
}

// CSRFMiddleware protects against CSRF attacks
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF check for safe methods
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Get session token
		sessionToken, err := getSessionToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get CSRF token from header or form
		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			csrfToken = r.FormValue("csrf_token")
		}

		if !ValidateCSRFToken(sessionToken, csrfToken) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetCSRFTokenHandler returns a new CSRF token for the session
func GetCSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := getSessionToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := GenerateCSRFToken(sessionToken)
	if err != nil {
		http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(token))
}

// CleanupExpiredCSRFTokens removes expired CSRF tokens
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			csrfStore.mu.Lock()
			now := time.Now()
			for session, token := range csrfStore.tokens {
				if now.After(token.ExpiresAt) {
					delete(csrfStore.tokens, session)
				}
			}
			csrfStore.mu.Unlock()
		}
	}()
}
