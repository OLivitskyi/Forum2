package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
)

// CreateCommentHandler handles comment creation.
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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
		PostID  string `json:"post_id"`
		Content string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to create comment with post ID: %s and content: %s", requestData.PostID, requestData.Content)

	if requestData.PostID == "" || requestData.Content == "" {
		http.Error(w, "Post ID and content are required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.FromString(requestData.PostID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = db.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Генерація comment_id
	commentID := uuid.Must(uuid.NewV4())
	log.Printf("Generated comment ID: %s", commentID)

	// Створення коментаря з використанням commentID
	err = db.CreateCommentWithID(commentID, postID, userID, requestData.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created comment with ID: %s", commentID)

	// Отримання новоствореного коментаря за його ID
	newComment, err := db.GetCommentByID(commentID)
	if err != nil {
		log.Printf("Error retrieving comment by ID %s: %v", commentID, err)
		http.Error(w, "Failed to retrieve comment after creation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newComment); err != nil {
		log.Printf("Failed to encode new comment to JSON: %v", err)
		http.Error(w, "Failed to encode comment", http.StatusInternalServerError)
	}
}

// GetCommentsHandler обробляє запити на отримання коментарів до конкретного посту
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Витягуємо ID посту з URL
	postIDStr := strings.TrimPrefix(r.URL.Path, "/api/post-comments/")
	if postIDStr == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		log.Printf("Invalid post ID format: %s", postIDStr)
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	comments, err := db.GetComments(postID)
	if err != nil {
		log.Printf("Failed to get comments for post ID %s: %v", postID, err)
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Printf("Failed to encode comments to JSON: %v", err)
		http.Error(w, "Failed to encode comments", http.StatusInternalServerError)
	}
}

// AddCommentReactionHandler handles adding reactions to comments.
func AddCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
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
		CommentID    string `json:"comment_id"`
		ReactionType string `json:"type"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.CommentID == "" || requestData.ReactionType == "" {
		http.Error(w, "Comment ID and reaction type are required", http.StatusBadRequest)
		return
	}

	commentID, err := uuid.FromString(requestData.CommentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	err = db.AddCommentReaction(userID, commentID, db.ReactionType(requestData.ReactionType))
	if err != nil {
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
