package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	// "log"

	// "strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const createtables string = `
CREATE TABLE IF NOT EXISTS categories (
	category_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
	category TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
	username TEXT NOT NULL UNIQUE ,
	age INTEGER ,
	gender TEXT ,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT DEFAULT NULL
);
CREATE TABLE IF NOT EXISTS posts (
	post_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
	user_id INTEGER NOT NULL ,
	title TEXT NOT NULL ,
	content TEXT NOT NULL ,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS comments (
	comment_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
	post_id INTEGER NOT NULL ,
	user_id INTEGER NOT NULL ,
	content TEXT NOT NULL ,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(post_id) REFERENCES posts(post_id),
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS likes (
	like_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL ,
    post_id INTEGER , 
    comment_id INTEGER , 
    user_id INTEGER NOT NULL , 
    type INTEGER NOT NULL , 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP , 
    FOREIGN KEY(post_id) REFERENCES posts(post_id) , 
    FOREIGN KEY(comment_id) REFERENCES comments(comment_id) , 
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS post_categories (
	post_id INTEGER NOT NULL, 
	category_id INTEGER NOT NULL, 
    FOREIGN KEY(post_id) REFERENCES posts(post_id) , 
    FOREIGN KEY(category_id) REFERENCES categories(category_id)
);
CREATE TABLE IF NOT EXISTS sessions (
	session_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	token TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	expires_at INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS messages (
	message_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	sender_id INTEGER NOT NULL,
	receiver_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	is_read BOOLEAN NOT NULL DEFAULT 0,
	FOREIGN KEY(sender_id) REFERENCES users(user_id),
	FOREIGN KEY(receiver_id) REFERENCES users(user_id)
);
CREATE TABLE IF NOT EXISTS user_status (
	user_id INTEGER PRIMARY KEY NOT NULL,
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

// func OpenDatabase() (*sql.DB, error) {
// 	dbPath := "./db/database.db"

// 	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {

// 		//Open a db connection
// 		db, err := sql.Open("sqlite3", dbPath)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// Tables
// 		if _, err := db.Exec(createtables); err != nil {
// 			return nil, err
// 		}

// 		return db, nil
// 	}
// 	db, err := sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

func RegisterUser(data []string) ([]User, error) {
	userList := []User{}

	if DB == nil {
		fmt.Println("can't connect to database")
		return nil, fmt.Errorf("db connection failed")
	}

	stmt, err := DB.Prepare(`INSERT INTO users (username, age, gender, firstname, lastname, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data[0], data[1], data[2], data[3], data[4], data[5], data[6])
	if err != nil {
		fmt.Println("error in stmt.exec, didn't write to the database because:")
		fmt.Println(err)
		return nil, err // potential security hole
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
		fmt.Printf("error finding username or email: %v\n", usernameOrEmail)
		return login, errors.New("can't find username or email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if err != nil {
		fmt.Println("Wrong password when logging in")
		return login, errors.New("wrong Password")

	}
	fmt.Println("login worked")
	return login, nil
}

func CreatePostDB(db *sql.DB, userID int, title, content string, categories string, createdAt time.Time) error {
	_, err := db.Exec("INSERT INTO posts (user_id, title, content, created_at) VALUES (?, ?, ?, ?)", userID, title, content, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func GetUserID(username string, db *sql.DB) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func AddMessage(senderID, receiverID int, content string) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}

	stmt, err := DB.Prepare(`INSERT INTO messages (sender_id, receiver_id, content, is_read) VALUES (?, ?, ?, 0)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(senderID, receiverID, content)
	if err != nil {
		return err
	}
	return nil
}

func GetMessages(senderID, receiverID int, limit, offset int) ([]Message, error) {
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

func UpdateUserStatus(userID int, isOnline bool) error {
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

func GetUserIDFromSession(token string) (int, error) {
	if DB == nil {
		return 0, fmt.Errorf("db connection failed")
	}

	var userID int
	err := DB.QueryRow(`SELECT user_id FROM sessions WHERE token = ?`, token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func MarkMessageAsRead(messageID int, userID int) error {
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
