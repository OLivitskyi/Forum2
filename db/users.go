package db

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser registers a new user in the database.
func RegisterUser(data []interface{}) error {
	if DB == nil {
		return fmt.Errorf("db connection failed")
	}
	stmt, err := DB.Prepare(`INSERT INTO users (user_id, username, age, gender, firstname, lastname, email, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Println("Prepare statement error:", err)
		return err
	}
	defer stmt.Close()
	userID, err := uuid.NewV4()
	if err != nil {
		log.Println("UUID generation error:", err)
		return err
	}
	data = append([]interface{}{userID.String()}, data...)
	_, err = stmt.Exec(data...)
	if err != nil {
		log.Println("Exec statement error:", err)
		return err
	}
	return nil
}

// LoginUser authenticates a user based on username or email and password.
func LoginUser(usernameOrEmail, password string) (Login, error) {
	var login Login
	var fieldname string
	if strings.Contains(usernameOrEmail, "@") {
		fieldname = "email"
	} else {
		fieldname = "username"
	}
	err := DB.QueryRow("SELECT username, email, password FROM users WHERE "+fieldname+" = ?", usernameOrEmail).Scan(&login.Username, &login.Email, &login.Password)
	if err != nil {
		return login, errors.New("can't find username or email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if err != nil {
		return login, errors.New("wrong Password")
	}
	return login, nil
}

// UpdateUserStatus updates the online status of a user.
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

// GetUserStatus retrieves the online status of all users.
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

// GetUserByID retrieves a user by their ID.
func GetUserByID(userID uuid.UUID) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT user_id, username, firstname, lastname, age, gender, email FROM users WHERE user_id = ?", userID).Scan(
		&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserIDByUsernameOrEmail retrieves a user ID based on their username or email.
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

// GetUsersOrderedByLastMessageOrAlphabetically retrieves users ordered by last message time or alphabetically.
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