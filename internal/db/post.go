package db

import (
	"database/sql"
	"fmt"
	"time"
)

func createPostsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "post" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"updatedAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"image"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}

func PostInsert(userID int, title, body string, categories []int) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("erreur lors du démarrage de la transaction: %v", err)
	}

	createSQL := `INSERT INTO post (user_id, title, body, status) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(createSQL, userID, title, body, "published")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("erreur lors de l'exécution de la requête: %v", err)
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("erreur lors de l'obtention de l'ID du dernier post inséré: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("erreur lors de l'engagement de la transaction: %v", err)
	}

	return int(postID), nil // Retourner l'ID du post et nil pour l'erreur
}

func PostUpdateContent(id int, body string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	updatedAt := time.Now()

	updateSQL := `UPDATE post SET body=?, updated_at=? WHERE id=?`
	_, err = tx.Exec(updateSQL, body, updatedAt, id)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

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
		tx.Rollback() // Rollback on error
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// CommentSelectByID retrieves a single comment for a specific comment ID from the database using a transaction
func PostTitleSelectById(postID int) (string, error) {
	db := SetupDatabase()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("error starting transaction: %v", err)
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
        SELECT p.title
        FROM post p
        WHERE p.id = ?`

	// Execute the query with the provided commentID
	var title string
	err = tx.QueryRow(query, postID).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no comment found with ID %d", postID)
		}
		return "", fmt.Errorf("error executing query: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error committing transaction: %v", err)
	}

	return title, nil
}

func UpdatePostStatus(id int, status string) error {
	// Setup the database connection
	db := SetupDatabase()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// SQL query to update only the status
	updateSQL := `UPDATE post SET status=? WHERE id=?`

	// Execute the update statement with the status and id parameters
	_, err = tx.Exec(updateSQL, status, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
