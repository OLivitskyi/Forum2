package handlers

import (
	"forum/db"
	"html/template"
	"log"
	"net/http"
	"time"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/registration" {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		MainPageHandler(w, r)
		return
	}
	if r.Method == "POST" {
		SignupProcess(w, r)
	}
}

func HomepageHandler(w http.ResponseWriter, r *http.Request) {
	if SessionExpired(w, r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.URL.Path != "/homepage" {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
	username := ValidateSession(w, r)
	if username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.Method == "GET" {
		MainPageHandler(w, r)
	}
}

func ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	username := ValidateSession(w, r)
	if username != "" && !SessionExpired(w, r) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
	user := ValidateSession(w, r)
	if user != "" {
		userID, err := db.GetUserIDByUsernameOrEmail(user)
		if err != nil {
			log.Printf("Failed to get user ID for logout: %v", err)
			http.Error(w, "Failed to logout", http.StatusInternalServerError)
			return
		}
		CloseSession(w, r)
		db.UpdateUserStatus(userID, false) // Update user status to offline
		broadcastUserStatus()
	}
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Logout successful"}`))

	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")
	login, err := db.LoginUser(usernameOrEmail, password)
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

	db.UpdateUserStatus(userID, true)
	broadcastUserStatus()
	log.Printf("User %s logged in with session token %s", login.Username, token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}
