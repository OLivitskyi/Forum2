package main

import (
	"fmt"
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	err := db.ConnectDatabase()
	if err != nil {
		fmt.Println("failed to connect to database in main.go")
		log.Fatal(err)
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Page routes
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/homepage", handlers.HomepageHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// API routes
	http.HandleFunc("/api/registration", handlers.SignupHandler)
	http.HandleFunc("/api/validate-session", handlers.ValidateSessionHandler)

	http.Handle("/api/create-category", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCategoryHandler)))
	http.Handle("/api/create-post", handlers.RequireLogin(http.HandlerFunc(handlers.CreatePostHandler)))
	http.Handle("/api/create-comment", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCommentHandler)))
	http.Handle("/api/get-posts", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostsHandler)))
	http.Handle("/api/get-comments", handlers.RequireLogin(http.HandlerFunc(handlers.GetCommentsHandler)))
	http.Handle("/api/send-message", handlers.RequireLogin(http.HandlerFunc(handlers.SendMessageHandler)))
	http.Handle("/api/get-messages", handlers.RequireLogin(http.HandlerFunc(handlers.GetMessagesHandler)))
	http.Handle("/api/update-status", handlers.RequireLogin(http.HandlerFunc(handlers.UpdateStatusHandler)))
	http.Handle("/api/get-user-status", handlers.RequireLogin(http.HandlerFunc(handlers.GetUserStatusHandler)))
	http.Handle("/api/mark-message-read", handlers.RequireLogin(http.HandlerFunc(handlers.MarkMessageAsReadHandler)))
	http.Handle("/api/get-users", handlers.RequireLogin(http.HandlerFunc(handlers.GetUsersHandler)))
	http.Handle("/api/add-post-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddPostReactionHandler)))
	http.Handle("/api/add-comment-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddCommentReactionHandler)))

	// WebSocket handler
	handlers.WebSocketHandler()

	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("Go to http://localhost:8080/\n")
	fmt.Printf("Ctrl + C to close the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
