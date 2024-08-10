package db

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
)

func CreateCommentWithID(commentID, postID, userID uuid.UUID, content string) error {
	_, err := DB.Exec("INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)", commentID, postID, userID, content, time.Now())
	return err
}

func GetComments(postID uuid.UUID) ([]Comment, error) {
	log.Printf("Fetching comments for post ID: %s", postID)
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
		log.Printf("Error querying comments: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		var username string
		err = rows.Scan(&c.ID, &c.UserID, &username, &c.Content, &c.CreatedAt, &c.LikeCount, &c.DislikeCount)
		if err != nil {
			log.Printf("Error scanning comment: %v", err)
			return nil, err
		}
		c.User = &User{
			Id:       c.UserID,
			Username: username,
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over comment rows: %v", err)
		return nil, err
	}

	log.Printf("Fetched %d comments for post ID: %s", len(comments), postID)
	return comments, nil
}

func AddCommentReaction(userID, commentID uuid.UUID, reactionType ReactionType) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, comment_id, type, created_at) VALUES (?, ?, ?, ?)", userID, commentID, reactionType, time.Now())
	return err
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

func GetCommentByID(commentID uuid.UUID) (*Comment, error) {
	log.Printf("Attempting to fetch comment with ID: %s", commentID)

	var comment Comment
	var username string
	err := DB.QueryRow(`
		SELECT c.comment_id, c.post_id, c.user_id, u.username, c.content, c.created_at,
		       COALESCE(SUM(CASE WHEN cr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
		       COALESCE(SUM(CASE WHEN cr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
		FROM comments c
		JOIN users u ON c.user_id = u.user_id
		LEFT JOIN likes cr ON c.comment_id = cr.comment_id
		WHERE c.comment_id = ?
		GROUP BY c.comment_id`, commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&username,
		&comment.Content,
		&comment.CreatedAt,
		&comment.LikeCount,
		&comment.DislikeCount,
	)

	if err != nil {
		log.Printf("Error fetching comment by ID %s: %v", commentID, err)
		return nil, err
	}

	comment.User = &User{
		Id:       comment.UserID,
		Username: username,
	}

	log.Printf("Successfully fetched comment with ID: %s", commentID)
	return &comment, nil
}
