package main

import (
	"fmt"
	"forum/db"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Підключення до бази даних
	err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database in main.go: %v", err)
	}

	// Очищення сесій
	err = db.ClearSessions()
	if err != nil {
		log.Fatalf("Failed to clear sessions: %v", err)
	}

	// Ініціалізація WebSocket
	handlers.InitWebSocketHandler()

	// Сервер статичних файлів
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршрути для сторінок
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/registration", handlers.SignupHandler)
	http.HandleFunc("/homepage", handlers.HomepageHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// API маршрути
	http.HandleFunc("/api/login", handlers.LoginHandler)
	http.HandleFunc("/api/validate-session", handlers.ValidateSessionHandler)
	http.Handle("/api/create-category", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCategoryHandler)))
	http.Handle("/api/get-categories", handlers.RequireLogin(http.HandlerFunc(handlers.GetCategoriesHandler)))
	http.Handle("/api/get-category", handlers.RequireLogin(http.HandlerFunc(handlers.GetCategoryByIDHandler)))
	http.Handle("/api/create-post", handlers.RequireLogin(http.HandlerFunc(handlers.CreatePostHandler)))

	// Маршрути для роботи з постами
	http.Handle("/api/posts", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostsHandler)))
	http.Handle("/api/post/", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostHandler)))

	// Маршрут для створення коментаря
	http.Handle("/api/post/comments/new", handlers.RequireLogin(http.HandlerFunc(handlers.CreateCommentHandler)))

	// Маршрут для отримання окремого посту
	http.Handle("/api/get-post/", handlers.RequireLogin(http.HandlerFunc(handlers.GetPostHandler)))

	// Маршрут для отримання коментарів до посту
	http.Handle("/api/post-comments/", handlers.RequireLogin(http.HandlerFunc(handlers.GetCommentsHandler)))

	// Інші API маршрути
	http.Handle("/api/send-message", handlers.RequireLogin(http.HandlerFunc(handlers.SendMessageHandler)))
	http.Handle("/api/get-messages", handlers.RequireLogin(http.HandlerFunc(handlers.GetMessagesHandler)))
	http.Handle("/api/update-status", handlers.RequireLogin(http.HandlerFunc(handlers.UpdateStatusHandler)))
	http.Handle("/api/get-user-status", handlers.RequireLogin(http.HandlerFunc(handlers.GetUserStatusHandler)))
	http.Handle("/api/mark-message-read", handlers.RequireLogin(http.HandlerFunc(handlers.MarkMessageAsReadHandler)))
	http.Handle("/api/get-users", handlers.RequireLogin(http.HandlerFunc(handlers.GetUsersHandler)))
	http.Handle("/api/add-comment-reaction", handlers.RequireLogin(http.HandlerFunc(handlers.AddCommentReactionHandler)))

	// WebSocket handler
	http.HandleFunc("/ws", handlers.WebSocketHandler)

	// Запуск сервера
	fmt.Printf("Starting server at port 8080\n")
	fmt.Printf("Go to http://localhost:8080/\n")
	fmt.Printf("Ctrl + C to close the server\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
