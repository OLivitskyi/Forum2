package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "Page not found.", http.StatusNotFound)
	// 	return
	// }
	if r.Method == "GET" {
		// fmt.Println("madis mainPageHandler " + r.URL.Path + " " + r.Method)
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
	// fmt.Println("magnus signUpHandler" + r.URL.Path + " " + r.Method)

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
	// fmt.Printf("this is the username we are using to signin:%v\n",username)
	if username == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// temporary (i think)
	if r.Method == "GET" {
		MainPageHandler(w, r)
		// ValidateSession(r)
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
		// need some work here still
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
