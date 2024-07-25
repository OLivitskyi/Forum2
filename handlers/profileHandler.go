package handlers

import "net/http"

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		http.Error(w, "Page not found.", http.StatusNotFound)
		return
	}
}
