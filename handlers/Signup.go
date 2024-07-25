package handlers

import (
	"fmt"
	"forum/db"
	"net/http"

	// "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	//"golang.org/x/crypto/bcrypt"
)

func SignupProcess(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("signupUsername")
	fmt.Println(username)
	email := r.FormValue("email")
	fmt.Println(email)
	firstName := r.FormValue("firstname")
	fmt.Println(firstName)
	lastName := r.FormValue("lastname")
	fmt.Println(lastName)
	age := r.FormValue("age")
	fmt.Println(age)
	gender := r.FormValue("gender")
	fmt.Println(gender)
	password := r.FormValue("signupPassword")
	fmt.Println(password)
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
	}
	fmt.Println(encryptedPassword)
	data := []string{username, age, gender, firstName, lastName, email, string(encryptedPassword)}
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
		// fmt.Fprintf(w, "Cant find username or email")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
	}
	// fmt.Printf("this is login:%v", login)

	// fmt.Println(username + "\n" + password)

	NewSession(w, login.Username)

	// rows, err := db.DB.Query("SELECT username, password, email FROM users WHERE username = ? OR email = ?")

}


// if its other text then internalservererror
