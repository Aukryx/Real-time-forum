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
	"user"	TEXT NOT NULL,
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

// Create - Insert a new comment
func CommentInsert(userID int, uuid string, postID int, body string) (*models.Comment, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	user := UserNicknameWithUUID(uuid)
	now := time.Now().Format("2006-01-02 15:04:05") // Fix date format

	// Match the column names in your 'comment' table
	insertSQL := `INSERT INTO comment (user_id, user, post_id, body, createdAt) 
                  VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(insertSQL, userID, user, postID, body, now)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error inserting comment: %v", err)
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	createdTime, _ := time.Parse("2006-01-02 15:04:05", now)
	comment := &models.Comment{
		ID:        int(commentID),
		UserID:    userID,
		Username:  user,
		PostID:    postID,
		Body:      body,
		CreatedAt: createdTime,
	}

	return comment, nil
}

// Read - Get comment by ID
func CommentSelectByID(commentID int) (*models.Comment, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, user, post_id, body, createdAt, updatedAt
             FROM comment WHERE id = ?`

	var comment models.Comment
	var createdAtStr, updatedAtStr string

	err = tx.QueryRow(query, commentID).Scan(
		&comment.ID, &comment.UserID, &comment.Username, &comment.PostID, &comment.Body,
		&createdAtStr, &updatedAtStr,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no comment found with ID %d", commentID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Parse time strings
	comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	comment.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &comment, nil
}

// Read - Get comments by post ID
func CommentSelectByPostID(postID int) ([]*models.Comment, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, user, post_id, body, createdAt, updatedAt
             FROM comment WHERE post_id = ? ORDER BY createdAt ASC`

	rows, err := tx.Query(query, postID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var createdAtStr, updatedAtStr string

		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.PostID, &comment.Body,
			&createdAtStr, &updatedAtStr); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}

		// Parse time strings
		comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		comment.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating comments: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return comments, nil
}

// Read - Get comments by user ID
func CommentSelectByUserID(userID int) ([]*models.Comment, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, user, post_id, body, createdAt, updatedAt
             FROM comment WHERE user_id = ? ORDER BY createdAt DESC`

	rows, err := tx.Query(query, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var createdAtStr, updatedAtStr string

		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Username, &comment.PostID, &comment.Body,
			&createdAtStr, &updatedAtStr); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}

		// Parse time strings
		comment.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		comment.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating comments: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return comments, nil
}

// Update - Update comment
func CommentUpdate(commentID int, body string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	updateSQL := `UPDATE comment SET body=?, updatedAt=? WHERE id=?`
	_, err = tx.Exec(updateSQL, body, now, commentID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete comment
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
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
