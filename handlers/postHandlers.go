package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

// CreatePostHandler handles post creation
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var requestData struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		CategoryIDs []string `json:"category_ids"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.Title == "" || requestData.Content == "" || len(requestData.CategoryIDs) == 0 {
		http.Error(w, "Title, content and at least one category are required", http.StatusBadRequest)
		return
	}

	var categoryIDInts []int
	for _, id := range requestData.CategoryIDs {
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		categoryIDInts = append(categoryIDInts, categoryID)
	}

	err = db.CreatePostDB(userID, requestData.Title, requestData.Content, categoryIDInts, time.Now())
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	post, err := db.GetPostByTitleAndContent(requestData.Title, requestData.Content)
	if err != nil {
		http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		return
	}

	postMessage := PostMessage{
		MessageID: uuid.Must(uuid.NewV4()),
		Post:      *post,
		Timestamp: time.Now(),
	}
	broadcast <- postMessage

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to get posts")

	posts, err := db.GetPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	// Перевірка на наявність постів
	if posts == nil {
		log.Println("No posts found")
		http.Error(w, "No posts found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully returned posts")
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Path[len("/api/get-post/"):]
	if postIDStr == "" {
		log.Println("Post ID is missing in the request")
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching post with ID: %s", postIDStr)
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		log.Printf("Invalid post ID format: %s", postIDStr)
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	post, err := db.GetPostByID(postID)
	if err != nil {
		log.Printf("Failed to get post: %v", err)
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully fetched post with ID: %s", postIDStr)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Printf("Error encoding post to JSON: %v", err)
		http.Error(w, "Failed to encode post to JSON", http.StatusInternalServerError)
		return
	}
}
