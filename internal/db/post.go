package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

func createPostsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "post" (
    "id"    INTEGER NOT NULL UNIQUE,
    "user_id"    INTEGER NOT NULL,
    "title"    TEXT NOT NULL,
    "body"    TEXT NOT NULL,
    "createdAt"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "image"    TEXT,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)

	// Also create post_category relation table if needed
	createPostCategoryTableSQL := `CREATE TABLE IF NOT EXISTS "post_category" (
	"post_id"	INTEGER NOT NULL,
	"category_id"	INTEGER NOT NULL,
	PRIMARY KEY("post_id", "category_id"),
	FOREIGN KEY("post_id") REFERENCES "post"("id") ON DELETE CASCADE,
	FOREIGN KEY("category_id") REFERENCES "category"("id")
)`
	executeSQL(db, createPostCategoryTableSQL)
}

type Post struct {
	ID        int
	UserID    int
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Image     string
}

// Create - Insert a new post
func PostInsert(userID int, title, body, imagePath string) (*models.Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	// Match the column names in your 'post' table
	insertSQL := `INSERT INTO post (user_id, title, body, image, createdAt) 
                  VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(insertSQL, userID, title, body, imagePath, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting post: %v", err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	post := &models.Post{
		ID:        int(postID),
		UserID:    userID,
		Title:     title,
		Body:      body,
		ImagePath: imagePath,
		CreatedAt: time.Now(),
	}

	return post, nil
}

// Read - Get post by ID
func PostSelectByID(postID int) (*Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, title, body, createdAt, updatedAt, image
             FROM post WHERE id = ?`

	var post Post
	var createdAtStr, updatedAtStr string

	err = tx.QueryRow(query, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Body,
		&createdAtStr, &updatedAtStr, &post.Image,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no post found with ID %d", postID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Parse time strings
	post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &post, nil
}

// Read - Get post title by ID (your existing function)
func PostTitleSelectById(postID int) (string, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT p.title FROM post p WHERE p.id = ?`

	var title string
	err = tx.QueryRow(query, postID).Scan(&title)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no post found with ID %d", postID)
		}
		return "", fmt.Errorf("error executing query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error committing transaction: %v", err)
	}

	return title, nil
}

// Read - Get all posts
func PostSelectAll() ([]models.Post, error) {
	db := SetupDatabase()
	defer db.Close()

	query := `SELECT id, user_id, title, body, image, createdAt FROM post ORDER BY createdAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var createdAtStr string // Declare as string first

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Body,
			&post.ImagePath,
			&post.CreatedAt, // This will now cause an error
		)
		// Modify the Scan to use a string
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Body,
			&post.ImagePath,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		// Parse the string to time.Time
		parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			// If parsing fails, use current time or handle as needed
			parsedTime = time.Now()
		}
		post.CreatedAt = parsedTime

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return posts, nil
}

// Read - Get posts by user ID
func PostSelectByUserID(userID int) ([]*Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, title, body, createdAt, updatedAt, image 
             FROM post WHERE user_id = ?`

	rows, err := tx.Query(query, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var createdAtStr, updatedAtStr string

		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body,
			&createdAtStr, &updatedAtStr, &post.Image); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		// Parse time strings
		post.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		post.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating posts: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return posts, nil
}

// Update - Update post content
func PostUpdateContent(id int, title, body string, image string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	updateSQL := `UPDATE post SET title=?, body=?, updatedAt=?, image=? WHERE id=?`
	_, err = tx.Exec(updateSQL, title, body, now, image, id)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete post
func PostDelete(postID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM post WHERE id=?`
	_, err = tx.Exec(deleteSQL, postID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
