package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

func createCommentsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "comment" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"updatedAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}

// CreateComment inserts a new comment into the database.
func CommentInsert(userID int, postID int, content string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Insert comment
	insertSQL := `INSERT INTO comment (post_id, user_id, content) VALUES (?, ?, ?)`
	result, err := tx.Exec(insertSQL, postID, userID, content)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	// Get the ID of the inserted comment
	commentID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Get the user ID of the post owner
	var postOwnerID int64
	err = tx.QueryRow("SELECT user_id FROM post WHERE id = ?", postID).Scan(&postOwnerID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting post owner ID: %v", err)
	}

	// Insert notification for the post owner
	_, err = tx.Exec("INSERT INTO notification (user_id, comment_id) VALUES (?, ?)", postOwnerID, commentID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting notification: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetCommentsByPostID retrieves all comments for a specific post ID from the database using a transaction
func CommentSelectByPostID(postID int, db *sql.DB) ([]models.Comment, error) {
	if db == nil {
		db := SetupDatabase()
		defer db.Close()
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure rollback in case of an error or panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rethrow panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
        FROM comment c
        JOIN user u ON c.user_id = u.id
        WHERE c.post_id = ?`

	// Execute the query with the provided postID
	rows, err := tx.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	for i, comment := range comments {
		comments[i].LikesDislikes, err = LikesSelectByCommentID(comment.ID, db)
		if err != nil {
			return nil, fmt.Errorf("error get likes on the comment: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return comments, nil
}

// UpdateComment updates an existing comment in the database.
func CommentUpdate(commentID int, userID int, postID int, content string) error {
	db := SetupDatabase()
	defer db.Close()

	updatedAt := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	updateSQL := `UPDATE comment SET content = ?, updated_at = ? WHERE user_id = ? AND post_id = ? AND id = ?`
	_, err = tx.Exec(updateSQL, content, updatedAt, userID, postID, commentID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func CommentDelete(commentID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM comment WHERE id=?`
	_, err = tx.Exec(deleteSQL, commentID)

	if err != nil {
		tx.Rollback() // Rollback on error
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// CommentSelectByID retrieves a single comment for a specific comment ID from the database using a transaction
func CommentSelectByID(commentID int64) (models.Comment, error) {
	db := SetupDatabase()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return models.Comment{}, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure rollback in case of an error or panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rethrow panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `
        SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
        FROM comment c
        JOIN user u ON c.user_id = u.id
        WHERE c.id = ?`

	// Execute the query with the provided commentID
	var comment models.Comment
	err = tx.QueryRow(query, commentID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Comment{}, fmt.Errorf("no comment found with ID %d", commentID)
		}
		return models.Comment{}, fmt.Errorf("error executing query: %v", err)
	}

	comment.LikesDislikes, err = LikesSelectByCommentID(comment.ID, db)
	if err != nil {
		return models.Comment{}, fmt.Errorf("error get likes on the comment: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.Comment{}, fmt.Errorf("error committing transaction: %v", err)
	}

	return comment, nil
}
