package handlers

import (
	"fmt"
	"forum/db"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignupProcess(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("signupUsername")
	email := r.FormValue("email")
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	age := r.FormValue("age")
	gender := r.FormValue("gender")
	password := r.FormValue("signupPassword")
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	data := []interface{}{username, age, gender, firstName, lastName, email, string(encryptedPassword)}
	_, err = db.RegisterUser(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf(w, "User created successfully!")
}

func LoginProcess(w http.ResponseWriter, r *http.Request) {
	usernameOrEmail := r.FormValue("username")
	password := r.FormValue("password")
	login, err := db.LoginUser(db.DB, usernameOrEmail, password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	userID, err := db.GetUserIDByUsernameOrEmail(usernameOrEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID"))
		return
	}
	token, err := NewSession(w, login.Username, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create session"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(24 * time.Hour),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
