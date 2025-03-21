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
    "user_id"    TEXT NOT NULL,
	"user"	TEXT NOT NULL,
    "title"    TEXT NOT NULL,
    "body"    TEXT NOT NULL,
    "createdAt"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"    DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}

// Create - Insert a new post
func PostInsert(userID int, uuid, title, body string) (*models.Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	user := UserNicknameWithUUID(uuid)
	now := time.Now().Format("2006-01-02 15:04:05") // Fix date format

	// Match the column names in your 'post' table
	insertSQL := `INSERT INTO post (user_id, user, title, body, createdAt) 
                  VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(insertSQL, userID, user, title, body, now)
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

	createdTime, _ := time.Parse("2006-01-02 15:04:05", now)
	post := &models.Post{
		ID:        int(postID),
		UserID:    userID,
		Username:  user,
		Title:     title,
		Body:      body,
		CreatedAt: createdTime,
	}

	return post, nil
}

// Read - Get post by ID
func PostSelectByID(postID int) (*models.Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, title, body, createdAt, updatedAt
             FROM post WHERE id = ?`

	var post models.Post
	var createdAtStr, updatedAtStr string

	err = tx.QueryRow(query, postID).Scan(
		&post.ID, &post.UserID, &post.Title, &post.Body,
		&createdAtStr, &updatedAtStr,
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

	query := `SELECT id, user_id, user, title, body, createdAt FROM post ORDER BY createdAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var createdAtStr string // Declare as string first

		// Modify the Scan to use a string
		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Username,
			&post.Title,
			&post.Body,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}

		// Parse the string to time.Time using correct format
		post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			// If parsing fails, use current time
			post.CreatedAt = time.Now()
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return posts, nil
}

// Read - Get posts by user ID
func PostSelectByUserID(userID int) ([]*models.Post, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, title, body, createdAt, updatedAt,
             FROM post WHERE user_id = ?`

	rows, err := tx.Query(query, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		var createdAtStr, updatedAtStr string

		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Body,
			&createdAtStr, &updatedAtStr); err != nil {
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
func PostUpdateContent(id int, title, body string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	updateSQL := `UPDATE post SET title=?, body=?, updatedAt=? WHERE id=?`
	_, err = tx.Exec(updateSQL, title, body, now, id)

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
