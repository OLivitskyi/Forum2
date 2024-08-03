package handlers

import (
	"encoding/json"
	"fmt"
	"forum/db"
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

	fmt.Println("Title:", requestData.Title)
	fmt.Println("Content:", requestData.Content)
	fmt.Println("Category IDs:", requestData.CategoryIDs)

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

	err = db.CreatePostDB(db.DB, userID, requestData.Title, requestData.Content, categoryIDInts, time.Now())
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// CreateCommentHandler handles comment creation
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

	if requestData.PostID == "" || requestData.Content == "" {
		http.Error(w, "Post ID and content are required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.FromString(requestData.PostID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = db.CreateComment(postID, userID, requestData.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetPostsHandler handles fetching posts
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GetCommentsHandler handles fetching comments for a post
func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	comments, err := db.GetComments(postID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// AddPostReactionHandler handles adding reactions to posts
func AddPostReactionHandler(w http.ResponseWriter, r *http.Request) {
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
		PostID       string `json:"post_id"`
		ReactionType string `json:"type"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.PostID == "" || requestData.ReactionType == "" {
		http.Error(w, "Post ID and reaction type are required", http.StatusBadRequest)
		return
	}

	postID, err := uuid.FromString(requestData.PostID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = db.AddPostReaction(userID, postID, db.ReactionType(requestData.ReactionType))
	if err != nil {
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddCommentReactionHandler handles adding reactions to comments
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
