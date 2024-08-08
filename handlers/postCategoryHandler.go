package handlers

import (
	"encoding/json"
	"forum/db"
	"net/http"
	"strconv"
)

// CreateCategoryHandler handles category creation
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	var requestData struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	if requestData.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}
	err = db.CreateCategory(requestData.Name)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetCategoryByIDHandler handles fetching a single category by ID
func GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Category ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}
	category, err := db.GetCategoryByID(id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// GetCategoriesHandler handles fetching all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	categories, err := db.GetCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
