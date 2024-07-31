package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const createtables string = `
CREATE TABLE IF NOT EXISTS categories (
	category_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	category TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY NOT NULL,
	username TEXT NOT NULL UNIQUE,
	age INTEGER,
	gender TEXT,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS posts (
	post_id UUID PRIMARY KEY NOT NULL,
	user_id UUID NOT NULL,
	subject TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS comments (
	comment_id UUID PRIMARY KEY NOT NULL,
	post_id UUID NOT NULL,
	user_id UUID NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(post_id) REFERENCES posts(post_id),
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS likes (
	like_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	post_id UUID,
	comment_id UUID,
	user_id UUID NOT NULL,
	type INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(post_id) REFERENCES posts(post_id),
	FOREIGN KEY(comment_id) REFERENCES comments(comment_id),
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS post_categories (
	post_id UUID NOT NULL,
	category_id INTEGER NOT NULL,
	FOREIGN KEY(post_id) REFERENCES posts(post_id),
	FOREIGN KEY(category_id) REFERENCES categories(category_id)
);
CREATE TABLE IF NOT EXISTS sessions (
	session_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	token TEXT NOT NULL,
	user_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	expires_at INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS messages (
	message_id UUID PRIMARY KEY NOT NULL,
	sender_id UUID NOT NULL,
	receiver_id UUID NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	is_read BOOLEAN NOT NULL DEFAULT 0,
	FOREIGN KEY(sender_id) REFERENCES users(user_id),
	FOREIGN KEY(receiver_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS user_status (
	user_id UUID PRIMARY KEY NOT NULL,
	is_online BOOLEAN NOT NULL DEFAULT 0,
	last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);`

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	_, err = DB.Exec(createtables)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func RegisterUser(data []interface{}) ([]User, error) {
	userList := []User{}
	if DB == nil {
		return nil, fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO users (user_id, username, age, gender, firstname, lastname, email, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return nil, err
	}
	defer stmt.Close()

	userID, err := uuid.NewV4()
	if err != nil {
		log.Println("UUID generation error:", err)
		return nil, err
	}

	data = append([]interface{}{userID.String()}, data...)

	_, err = stmt.Exec(data...)
	if err != nil {
		log.Println("Exec statement error:", err)
		return nil, err
	}
	return userList, nil
}

func LoginUser(db *sql.DB, usernameOrEmail, password string) (Login, error) {
	var login Login
	var fieldname string
	if strings.Contains(usernameOrEmail, "@") {
		fieldname = "email"
	} else {
		fieldname = "username"
	}
	err := db.QueryRow("SELECT username, email, password FROM users WHERE "+fieldname+" = ?", usernameOrEmail).Scan(&login.Username, &login.Email, &login.Password)
	if err != nil {
		return login, errors.New("can't find username or email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if err != nil {
		return login, errors.New("wrong Password")
	}
	return login, nil
}

func CreatePostDB(db *sql.DB, userID uuid.UUID, subject, content string, categoryIDs []int, createdAt time.Time) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	postID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO posts (post_id, user_id, subject, content, created_at) VALUES (?, ?, ?, ?, ?)", postID, userID, subject, content, createdAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, categoryID := range categoryIDs {
		_, err = tx.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetPosts() ([]Post, error) {
	rows, err := DB.Query(`
        SELECT p.post_id, p.user_id, p.subject, p.content, p.created_at, 
               COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
               COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
        FROM posts p
        LEFT JOIN likes pr ON p.post_id = pr.post_id
        GROUP BY p.post_id
        ORDER BY p.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Content, &p.CreatedAt, &p.LikeCount, &p.DislikeCount)
		if err != nil {
			return nil, err
		}
		user, err := GetUserByID(p.UserID)
		if err != nil {
			return nil, err
		}
		p.User = user

		categories, err := GetPostCategories(p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = convertToCategoryPointers(categories)

		comments, err := GetComments(p.ID)
		if err != nil {
			return nil, err
		}
		p.Comments = convertToCommentPointers(comments)

		posts = append(posts, p)
	}
	return posts, nil
}

func convertToCategoryPointers(categories []Category) []*Category {
	var categoryPointers []*Category
	for i := range categories {
		categoryPointers = append(categoryPointers, &categories[i])
	}
	return categoryPointers
}

func convertToCommentPointers(comments []Comment) []*Comment {
	var commentPointers []*Comment
	for i := range comments {
		commentPointers = append(commentPointers, &comments[i])
	}
	return commentPointers
}

func GetPostCategories(postID uuid.UUID) ([]Category, error) {
	rows, err := DB.Query(`
        SELECT c.category_id, c.category
        FROM categories c
        JOIN post_categories pc ON c.category_id = pc.category_id
        WHERE pc.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func CreateComment(postID, userID uuid.UUID, content string) error {
	_, err := DB.Exec("INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)", uuid.Must(uuid.NewV4()), postID, userID, content, time.Now())
	return err
}

func GetComments(postID uuid.UUID) ([]Comment, error) {
	rows, err := DB.Query(`
        SELECT c.comment_id, c.user_id, u.username, c.content, c.created_at,
               COALESCE(SUM(CASE WHEN cr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
               COALESCE(SUM(CASE WHEN cr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
        LEFT JOIN likes cr ON c.comment_id = cr.comment_id
        WHERE c.post_id = ?
        GROUP BY c.comment_id
        ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.UserID, &c.User.Username, &c.Content, &c.CreatedAt, &c.LikeCount, &c.DislikeCount)
		if err != nil {
			return nil, err
		}
		user, err := GetUserByID(c.UserID)
		if err != nil {
			return nil, err
		}
		c.User = user

		comments = append(comments, c)
	}
	return comments, nil
}

func GetUserByID(userID uuid.UUID) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT user_id, username, firstname, lastname, age, gender, email FROM users WHERE user_id = ?", userID).Scan(
		&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateCategory(name string) error {
	_, err := DB.Exec("INSERT INTO categories (category) VALUES (?)", name)
	return err
}

func GetCategories() ([]Category, error) {
	rows, err := DB.Query("SELECT category_id, category FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func AddMessage(senderID, receiverID uuid.UUID, content string) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO messages (message_id, sender_id, receiver_id, content, created_at, is_read) VALUES (?, ?, ?, ?, ?, 0)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	messageID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(messageID, senderID, receiverID, content, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func GetMessages(senderID, receiverID uuid.UUID, limit, offset int) ([]Message, error) {
	if DB == nil {
		return nil, fmt.Errorf("db connection failed")
	}
	rows, err := DB.Query(`SELECT message_id, sender_id, receiver_id, content, created_at, is_read FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) ORDER BY created_at DESC LIMIT ? OFFSET ?`, senderID, receiverID, receiverID, senderID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.IsRead)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func UpdateUserStatus(userID uuid.UUID, isOnline bool) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO user_status (user_id, is_online) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET is_online = ?, last_activity = CURRENT_TIMESTAMP`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userID, isOnline, isOnline)
	if err != nil {
		return err
	}
	return nil
}

func GetUserStatus() ([]UserStatus, error) {
	if DB == nil {
		return nil, fmt.Errorf("db connection failed")
	}
	rows, err := DB.Query(`SELECT user_id, is_online, last_activity FROM user_status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var statuses []UserStatus
	for rows.Next() {
		var status UserStatus
		err := rows.Scan(&status.UserID, &status.IsOnline, &status.LastActivity)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func GetUsersOrderedByLastMessageOrAlphabetically() ([]User, error) {
	if DB == nil {
		return nil, fmt.Errorf("db connection failed")
	}
	rows, err := DB.Query(`
	SELECT users.user_id, users.username, MAX(messages.created_at) AS last_message_time
	FROM users
	LEFT JOIN messages ON users.user_id = messages.sender_id OR users.user_id = messages.receiver_id
	GROUP BY users.user_id
	ORDER BY last_message_time DESC, users.username ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserIDFromSession(token string) (uuid.UUID, error) {
	if DB == nil {
		return uuid.Nil, fmt.Errorf("db connection failed")
	}
	var userID uuid.UUID
	err := DB.QueryRow(`SELECT user_id FROM sessions WHERE token = ?`, token).Scan(&userID)
	if err != nil {
		log.Printf("Error finding session token: %v", err)
		return uuid.Nil, err
	}
	log.Printf("Session token found for user ID: %d", userID)
	return userID, nil
}

func MarkMessageAsRead(messageID uuid.UUID, userID uuid.UUID) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`UPDATE messages SET is_read = 1 WHERE message_id = ? AND receiver_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(messageID, userID)
	if err != nil {
		return err
	}
	return nil
}

func SaveSession(token string, userID uuid.UUID, expiration time.Time) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token, userID, expiration.Unix())
	if err != nil {
		return err
	}
	return nil
}

func DeleteSession(token string) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`DELETE FROM sessions WHERE token = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(token)
	if err != nil {
		return err
	}
	return nil
}

func GetUserID(usernameOrEmail string, db *sql.DB) (uuid.UUID, error) {
	var userID uuid.UUID
	var fieldname string
	if strings.Contains(usernameOrEmail, "@") {
		fieldname = "email"
	} else {
		fieldname = "username"
	}
	err := db.QueryRow("SELECT user_id FROM users WHERE "+fieldname+" = ?", usernameOrEmail).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

// add reaction from socket
func AddPostReaction(userID, postID uuid.UUID, reactionType ReactionType) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, post_id, type, created_at) VALUES (?, ?, ?, ?)", userID, postID, reactionType, time.Now())
	return err
}

func AddCommentReaction(userID, commentID uuid.UUID, reactionType ReactionType) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, comment_id, type, created_at) VALUES (?, ?, ?, ?)", userID, commentID, reactionType, time.Now())
	return err
}

func GetPostReactions(postID uuid.UUID) ([]Reaction, error) {
	rows, err := DB.Query("SELECT user_id, post_id, type FROM likes WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var r Reaction
		err := rows.Scan(&r.UserID, &r.PostID, &r.Type)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}
	return reactions, nil
}

func GetCommentReactions(commentID uuid.UUID) ([]Reaction, error) {
	rows, err := DB.Query("SELECT user_id, comment_id, type FROM likes WHERE comment_id = ?", commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []Reaction
	for rows.Next() {
		var r Reaction
		err := rows.Scan(&r.UserID, &r.CommentID, &r.Type)
		if err != nil {
			return nil, err
		}
		reactions = append(reactions, r)
	}
	return reactions, nil
}

func GetUserIDByUsernameOrEmail(usernameOrEmail string) (uuid.UUID, error) {
	var userID uuid.UUID
	var fieldname string
	if strings.Contains(usernameOrEmail, "@") {
		fieldname = "email"
	} else {
		fieldname = "username"
	}
	err := DB.QueryRow("SELECT user_id FROM users WHERE "+fieldname+" = ?", usernameOrEmail).Scan(&userID)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
