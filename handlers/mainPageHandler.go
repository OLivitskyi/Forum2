package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
	}
	if r.Method == "POST" {
		LoginProcess(w, r)
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
	if SessionExpired(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.URL.Path != "/homepage" {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
	username := ValidateSession(r)
	if username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		MainPageHandler(w, r)
	}
}

func ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	username := ValidateSession(r)
	if username != "" && !SessionExpired(r) {
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

	user := ValidateSession(r)
	if user != "" {
		CloseSession(w, r)
	}

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Logout successful"}`))
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
