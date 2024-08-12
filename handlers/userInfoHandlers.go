package handlers

import (
	"encoding/json"
	"forum/db"
	"log"
	"net/http"
)

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Отримати інформацію про користувача
	user, err := db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Створити відповідь із потрібною інформацією
	userInfo := map[string]interface{}{
		"user_id":   user.Id,
		"username":  user.Username,
		"firstname": user.FirstName,
		"lastname":  user.LastName,
	}

	// Відправити відповідь у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
