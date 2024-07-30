package handlers

import (
	"fmt"
	"forum/db"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func SignupProcess(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("signupUsername")
	email := r.FormValue("email")
	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	ageStr := r.FormValue("age")
	gender := r.FormValue("gender")
	password := r.FormValue("signupPassword")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Invalid age format", http.StatusBadRequest)
		return
	}

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
	username := r.FormValue("username")
	password := r.FormValue("password")
	login, err := db.LoginUser(db.DB, username, password)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	userID, err := db.GetUserID(username, db.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get user ID"))
		return
	}

	NewSession(w, login.Username, userID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
