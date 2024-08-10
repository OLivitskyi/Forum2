package main

import (
	"fmt"
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Connect to the database
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database in main.go: %v", err)
	}

	// Clear sessions and user statuses
	if err := db.ClearSessions(); err != nil {
		log.Fatalf("Failed to clear sessions: %v", err)
	}

	if err := db.ClearUserStatus(); err != nil {
		log.Fatalf("Failed to clear user_status: %v", err)
	}

	// Initialize WebSocket handler
	handlers.InitWebSocketHandler()

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Page routes
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/registration", handlers.SignupHandler)
	http.HandleFunc("/homepage", handlers.HomepageHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// API routes
	apiRoutes()

	// WebSocket handler
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	// Start the server
	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("Go to http://localhost:8080/\n")
	fmt.Printf("Ctrl + C to close the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func apiRoutes() {
	// Authentication routes
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/validate-session", handlers.ValidateSessionHandler)

	// Category routes
	http.Handle("/api/create-category", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCategoryHandler)))
	http.Handle("/api/get-categories", handlers.RequireLogin(http.HandlerFunc(handlers.GetCategoriesHandler)))
	http.Handle("/api/get-category", handlers.RequireLogin(http.HandlerFunc(handlers.GetCategoryByIDHandler)))

	// Post routes
	http.Handle("/api/create-post", handlers.RequireLogin(http.HandlerFunc(handlers.CreatePostHandler)))
	http.Handle("/api/posts", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostsHandler)))
	http.Handle("/api/post/", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostHandler)))
	http.Handle("/api/get-post/", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostHandler)))

	// Comment routes
	http.Handle("/api/post/comments/new", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCommentHandler)))
	http.Handle("/api/post-comments/", handlers.RequireLogin(http.HandlerFunc(handlers.GetCommentsHandler)))
	http.Handle("/api/add-comment-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddCommentReactionHandler)))

	// Message routes
	http.Handle("/api/send-message", handlers.RequireLogin(http.HandlerFunc(handlers.SendMessageHandler)))
	http.Handle("/api/get-messages", handlers.RequireLogin(http.HandlerFunc(handlers.GetMessagesHandler)))
	http.Handle("/api/mark-message-read", handlers.RequireLogin(http.HandlerFunc(handlers.MarkMessageAsReadHandler)))

	// User status routes
	http.Handle("/api/update-status", handlers.RequireLogin(http.HandlerFunc(handlers.UpdateStatusHandler)))
	http.Handle("/api/get-user-status", handlers.RequireLogin(http.HandlerFunc(handlers.GetUserStatusHandler)))
	http.Handle("/api/get-users", handlers.RequireLogin(http.HandlerFunc(handlers.GetUsersHandler)))
}
