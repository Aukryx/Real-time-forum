package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

func createLikesDislikesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "like_dislike" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER,
	"comment_id"	INTEGER,
	"status"	INTEGER NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment_id") REFERENCES "comment"("id"),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}

func LikeDislikeInsert(userID int, postID *int, commentID *int, status int) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	createSQL := `INSERT INTO like_dislike (user_id, post_id, comment_id, status) 
                  VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(createSQL, userID, postID, commentID, status)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	likeDislikeID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error getting last inserted like/dislike ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return int(likeDislikeID), nil
}

// Read - Get a like/dislike by ID
func LikeDislikeSelectByID(likeDislikeID int) (*models.LikeDislike, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, post_id, comment_id, status, createdAt
              FROM like_dislike WHERE id = ?`

	var likeDislike models.LikeDislike
	var createdAtStr string
	var postID, commentID sql.NullInt64

	err = tx.QueryRow(query, likeDislikeID).Scan(
		&likeDislike.ID, &likeDislike.UserID, &postID, &commentID,
		&likeDislike.Status, &createdAtStr,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no like/dislike found with ID %d", likeDislikeID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Handle nullable fields
	if postID.Valid {
		likeDislike.PostID = int(postID.Int64)
	}
	if commentID.Valid {
		likeDislike.CommentID = int(commentID.Int64)
	}

	// Parse time string
	likeDislike.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &likeDislike, nil
}

// Read - Get likes/dislikes by post ID
func LikeDislikeSelectByPostID(postID int) ([]*models.LikeDislike, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, post_id, comment_id, status, createdAt
              FROM like_dislike WHERE post_id = ?`

	rows, err := tx.Query(query, postID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var likeDislikes []*models.LikeDislike
	for rows.Next() {
		likeDislike := &models.LikeDislike{}
		var createdAtStr string
		var commentID sql.NullInt64

		if err := rows.Scan(&likeDislike.ID, &likeDislike.UserID, &likeDislike.PostID,
			&commentID, &likeDislike.Status, &createdAtStr); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning like/dislike: %v", err)
		}

		// Handle nullable fields
		if commentID.Valid {
			likeDislike.CommentID = int(commentID.Int64)
		}

		// Parse time string
		likeDislike.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

		likeDislikes = append(likeDislikes, likeDislike)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating like/dislikes: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return likeDislikes, nil
}

// Update - Update like/dislike status
func LikeDislikeUpdateStatus(likeDislikeID int, status int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	updateSQL := `UPDATE like_dislike SET status = ? WHERE id = ?`
	_, err = tx.Exec(updateSQL, status, likeDislikeID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete like/dislike
func LikeDislikeDelete(likeDislikeID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM like_dislike WHERE id = ?`
	_, err = tx.Exec(deleteSQL, likeDislikeID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
