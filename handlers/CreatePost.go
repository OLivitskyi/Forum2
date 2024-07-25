package handlers

import (
	"forum/db"
	"log"
	"net/http"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	postTitle := r.FormValue("postTitle")
	postContent := r.FormValue("postText")
	createdAt := time.Now()
	author := r.FormValue("author")
	postCategories := r.FormValue("categories")
	var userID int
	userID, err := db.GetUserID(author, db.DB)
	if err != nil {
		log.Println("Database error:", err)
	}
	err = db.CreatePostDB(db.DB, userID, postTitle, postContent,postCategories, createdAt)
	if err != nil {
		log.Println("error creating post")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
	}
}
