package db

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
)

// CreatePostDB creates a new post in the database along with its categories.
func CreatePostDB(userID uuid.UUID, subject, content string, categoryIDs []int, createdAt time.Time) error {
	tx, err := DB.Begin()
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

// GetPosts retrieves all posts from the database.
func GetPosts() ([]Post, error) {
	log.Println("Fetching all posts")
	rows, err := DB.Query(`
        SELECT p.post_id, p.user_id, p.subject, p.content, p.created_at, 
               COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
               COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
        FROM posts p
        LEFT JOIN likes pr ON p.post_id = pr.post_id
        GROUP BY p.post_id
        ORDER BY p.created_at DESC`)
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		return nil, err
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Subject, &p.Content, &p.CreatedAt, &p.LikeCount, &p.DislikeCount)
		if err != nil {
			log.Printf("Error scanning post: %v", err)
			return nil, err
		}
		user, err := GetUserByID(p.UserID)
		if err != nil {
			log.Printf("Error getting user by ID: %v", err)
			return nil, err
		}
		p.User = user
		categories, err := GetPostCategories(p.ID)
		if err != nil {
			log.Printf("Error getting post categories: %v", err)
			return nil, err
		}
		p.Categories = convertToCategoryPointers(categories)
		comments, err := GetComments(p.ID)
		if err != nil {
			log.Printf("Error getting comments for post: %v", err)
			return nil, err
		}
		p.Comments = convertToCommentPointers(comments)
		posts = append(posts, p)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over post rows: %v", err)
		return nil, err
	}
	log.Printf("Fetched %d posts", len(posts))
	return posts, nil
}

// GetPostByID retrieves a post by its ID.
func GetPostByID(postID uuid.UUID) (*Post, error) {
	log.Printf("Fetching post by ID: %s", postID)
	var post Post
	err := DB.QueryRow(`
        SELECT p.post_id, p.user_id, p.subject, p.content, p.created_at, 
               COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
               COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
        FROM posts p
        LEFT JOIN likes pr ON p.post_id = pr.post_id
        WHERE p.post_id = ?
        GROUP BY p.post_id
        ORDER BY p.created_at DESC`, postID).Scan(&post.ID, &post.UserID, &post.Subject, &post.Content, &post.CreatedAt, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		log.Printf("Error getting post by ID: %v", err)
		return nil, err
	}

	user, err := GetUserByID(post.UserID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}
	post.User = user

	categories, err := GetPostCategories(post.ID)
	if err != nil {
		log.Printf("Error getting post categories: %v", err)
		return nil, err
	}
	post.Categories = convertToCategoryPointers(categories)

	comments, err := GetComments(post.ID)
	if err != nil {
		log.Printf("Error getting comments for post: %v", err)
		return nil, err
	}
	post.Comments = convertToCommentPointers(comments)

	log.Printf("Fetched post: %+v", post)
	return &post, nil
}

// GetPostByTitleAndContent retrieves a post by its title and content.
func GetPostByTitleAndContent(title, content string) (*Post, error) {
	var post Post
	err := DB.QueryRow(`
        SELECT p.post_id, p.user_id, p.subject, p.content, p.created_at, 
               COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0) AS like_count,
               COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislike_count
        FROM posts p
        LEFT JOIN likes pr ON p.post_id = pr.post_id
        WHERE p.subject = ? AND p.content = ?
        GROUP BY p.post_id
        ORDER BY p.created_at DESC`, title, content).Scan(&post.ID, &post.UserID, &post.Subject, &post.Content, &post.CreatedAt, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		return nil, err
	}

	user, err := GetUserByID(post.UserID)
	if err != nil {
		return nil, err
	}
	post.User = user

	categories, err := GetPostCategories(post.ID)
	if err != nil {
		return nil, err
	}
	post.Categories = convertToCategoryPointers(categories)

	comments, err := GetComments(post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = convertToCommentPointers(comments)

	return &post, nil
}

// AddPostReaction adds a reaction to a post.
func AddPostReaction(userID, postID uuid.UUID, reactionType ReactionType) error {
	_, err := DB.Exec("INSERT INTO likes (user_id, post_id, type, created_at) VALUES (?, ?, ?, ?)", userID, postID, reactionType, time.Now())
	return err
}

// GetPostReactions retrieves reactions for a given post.
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

// GetPostCategories retrieves categories associated with a given post.
func GetPostCategories(postID uuid.UUID) ([]Category, error) {
	log.Printf("Fetching categories for post ID: %s", postID)
	rows, err := DB.Query(`
        SELECT c.category_id, c.category
        FROM categories c
        JOIN post_categories pc ON c.category_id = pc.category_id
        WHERE pc.post_id = ?`, postID)
	if err != nil {
		log.Printf("Error querying post categories: %v", err)
		return nil, err
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			log.Printf("Error scanning category: %v", err)
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over category rows: %v", err)
		return nil, err
	}
	log.Printf("Fetched %d categories for post ID: %s", len(categories), postID)
	return categories, nil
}

// Helper functions to convert slices to pointer slices.
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
