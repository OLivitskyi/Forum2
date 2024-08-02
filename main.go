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
	http.HandleFunc("/send-message", handlers.SendMessageHandler)
	http.HandleFunc("/get-messages", handlers.GetMessagesHandler)
	http.HandleFunc("/update-status", handlers.UpdateStatusHandler)
	http.HandleFunc("/get-user-status", handlers.GetUserStatusHandler)
	http.HandleFunc("/mark-message-read", handlers.MarkMessageAsReadHandler)
	http.HandleFunc("/get-users", handlers.GetUsersHandler)

	// posts and comments
	http.HandleFunc("/create-category", handlers.CreateCategoryHandler)
	http.Handle("/create-post", handlers.RequireLogin(http.HandlerFunc(handlers.CreatePostHandler)))
	http.HandleFunc("/create-comment", handlers.CreateCommentHandler)
	http.HandleFunc("/get-posts", handlers.GetPostsHandler)
	http.HandleFunc("/get-comments", handlers.GetCommentsHandler)

	// reactions
	http.HandleFunc("/add-post-reaction", handlers.AddPostReactionHandler)
	http.HandleFunc("/add-comment-reaction", handlers.AddCommentReactionHandler)

	handlers.WebSocketHandler()

	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("Go to http://localhost:8080/\n")
	fmt.Printf("Ctrl + C to close the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
