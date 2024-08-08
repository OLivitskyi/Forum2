package db

import (
	"time"

	"github.com/gofrs/uuid"
)

// CreateComment creates a new comment for a post.
func CreateComment(postID, userID uuid.UUID, content string) error {
	_, err := DB.Exec("INSERT INTO comments (comment_id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)", uuid.Must(uuid.NewV4()), postID, userID, content, time.Now())
	return err
}

// GetComments retrieves all comments for a given post.
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

// AddCommentReaction adds a reaction to a comment.
func AddCommentReaction(userID, commentID uuid.UUID, reactionType ReactionType) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, comment_id, type, created_at) VALUES (?, ?, ?, ?)", userID, commentID, reactionType, time.Now())
	return err
}

// GetCommentReactions retrieves reactions for a given comment.
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
