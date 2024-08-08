package handlers

import (
	"forum/db"
	"log"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")
	login, err := db.LoginUser(db.DB, usernameOrEmail, password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		log.Printf("Login failed: %v", err)
		return
	}
	userID, err := db.GetUserIDByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID"))
		log.Printf("Failed to get user ID: %v", err)
		return
	}
	token, err := NewSession(w, login.Username, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create session"))
		log.Printf("Failed to create session: %v", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	})

	log.Printf("User %s logged in with session token %s", login.Username, token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
