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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/registration", handlers.SignupHandler)
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/homepage", handlers.HomepageHandler)
	http.HandleFunc("/validate-session", handlers.ValidateSessionHandler)

	// messages
	http.HandleFunc("/api/send-message", handlers.SendMessageHandler)
	http.HandleFunc("/api/get-messages", handlers.GetMessagesHandler)
	http.HandleFunc("/api/update-status", handlers.UpdateStatusHandler)
	http.HandleFunc("/api/get-user-status", handlers.GetUserStatusHandler)
	http.HandleFunc("/api/mark-message-read", handlers.MarkMessageAsReadHandler)
	http.HandleFunc("/api/get-users", handlers.GetUsersHandler)

	// posts and comments
	http.Handle("/api/create-category", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCategoryHandler)))
	http.Handle("/api/create-post", handlers.RequireLogin(http.HandlerFunc(handlers.CreatePostHandler)))
	http.Handle("/api/create-comment", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCommentHandler)))
	http.Handle("/api/get-posts", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostsHandler)))
	http.Handle("/api/get-comments", handlers.RequireLogin(http.HandlerFunc(handlers.GetCommentsHandler)))

	// reactions
	http.Handle("/api/add-post-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddPostReactionHandler)))
	http.Handle("/api/add-comment-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddCommentReactionHandler)))

	handlers.WebSocketHandler()

	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("Go to http://localhost:8080/\n")
	fmt.Printf("Ctrl + C to close the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
