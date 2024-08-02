package handlers

import (
	"encoding/json"
	"forum/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

// CreateCategoryHandler handles category creation
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	categoryName := r.FormValue("name")
	if categoryName == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}
	err := db.CreateCategory(categoryName)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// CreatePostHandler handles post creation
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryIDs := r.Form["category_ids"]
	if title == "" || content == "" || len(categoryIDs) == 0 {
		http.Error(w, "Title, content and at least one category are required", http.StatusBadRequest)
		return
	}
	var categoryIDInts []int
	for _, id := range categoryIDs {
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		categoryIDInts = append(categoryIDInts, categoryID)
	}
	err = db.CreatePostDB(db.DB, userID, title, content, categoryIDInts, time.Now())
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
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")
	if postIDStr == "" || content == "" {
		http.Error(w, "Post ID and content are required", http.StatusBadRequest)
		return
	}
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	err = db.CreateComment(postID, userID, content)
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
	postIDStr := r.FormValue("post_id")
	reactionType := r.FormValue("type")
	if postIDStr == "" || reactionType == "" {
		http.Error(w, "Post ID and reaction type are required", http.StatusBadRequest)
		return
	}
	postID, err := uuid.FromString(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	err = db.AddPostReaction(userID, postID, db.ReactionType(reactionType))
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
	commentIDStr := r.FormValue("comment_id")
	reactionType := r.FormValue("type")
	if commentIDStr == "" || reactionType == "" {
		http.Error(w, "Comment ID and reaction type are required", http.StatusBadRequest)
		return
	}
	commentID, err := uuid.FromString(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	err = db.AddCommentReaction(userID, commentID, db.ReactionType(reactionType))
	if err != nil {
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
