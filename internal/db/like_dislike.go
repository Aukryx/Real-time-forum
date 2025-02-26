package db

import (
	"database/sql"
	"fmt"
	"models"
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

// Create likes and dislikes in the database.
func LikesInsert(userID, postID, commentID int, isLike bool) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	var stmt *sql.Stmt
	var result sql.Result

	if postID != -1 {
		// Dealing with post
		stmt, err = tx.Prepare(`
            INSERT INTO like_dislike (user_id, post_id, is_like)
            VALUES (?, ?, ?)
            ON CONFLICT(user_id, post_id) DO UPDATE SET is_like = ?
        `)
		if err != nil {
			return fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()

		result, err = stmt.Exec(userID, postID, isLike, isLike)
	} else if commentID != -1 {
		// Dealing with comment
		stmt, err = tx.Prepare(`
            INSERT INTO like_dislike (user_id, comment_id, is_like)
            VALUES (?, ?, ?)
            ON CONFLICT(user_id, comment_id) DO UPDATE SET is_like = ?
        `)
		if err != nil {
			return fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()

		result, err = stmt.Exec(userID, commentID, isLike, isLike)
	} else {
		return fmt.Errorf("both postID and commentID cannot be -1")
	}

	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("no rows were affected")
	}
	if postID != -1 {
		// Get the ID of the inserted comment
		var likeID int64
		err = tx.QueryRow("SELECT id FROM like_dislike WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeID)
		if err != nil {
			return fmt.Errorf("error getting like/dislike ID: %v", err)
		}

		// Get the user ID of the post owner
		var postOwnerID int64
		err = tx.QueryRow("SELECT user_id FROM post WHERE id = ?", postID).Scan(&postOwnerID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error getting post owner ID: %v", err)
		}

		// Insert notification for the post owner
		_, err = tx.Exec("INSERT INTO notifications (user_id, like_dislike_id) VALUES (?, ?)", postOwnerID, likeID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting notification: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func LikesSelectByPostID(postID int, db *sql.DB) ([]models.LikesDislikes, error) {
	if db == nil {
		db := SetupDatabase()
		defer db.Close()
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure rollback in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rethrow panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT *
              FROM like_dislike
              WHERE post_id=?`

	rows, err := tx.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var likesDislikes []models.LikesDislikes
	for rows.Next() {
		var likeDislike models.LikesDislikes
		if err := rows.Scan(&likeDislike.ID, &likeDislike.UserID, &likeDislike.PostID, &likeDislike.CommentID, &likeDislike.IsLike, &likeDislike.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		likesDislikes = append(likesDislikes, likeDislike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return likesDislikes, nil
}

func LikesSelectByCommentID(commentID int, db *sql.DB) ([]models.LikesDislikes, error) {
	if db == nil {
		db := SetupDatabase()
		defer db.Close()
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure rollback in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rethrow panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT *
              FROM like_dislike
              WHERE comment_id=?`

	rows, err := tx.Query(query, commentID)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var likesDislikes []models.LikesDislikes
	for rows.Next() {
		var likeDislike models.LikesDislikes
		if err := rows.Scan(&likeDislike.ID, &likeDislike.UserID, &likeDislike.PostID, &likeDislike.CommentID, &likeDislike.IsLike, &likeDislike.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		likesDislikes = append(likesDislikes, likeDislike)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return likesDislikes, nil
}

// Create likes and dislikes in the database.
func LikesUpdate(userID, postID, commentID int, isLike bool) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	if postID == -1 {
		insertSQL := `UPDATE like_dislike SET is_like=? WHERE user_id=? AND comment_id=?`
		_, err = tx.Exec(insertSQL, isLike, userID, commentID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing statement: %v", err)
		}

	} else if commentID == -1 {
		insertSQL := `UPDATE like_dislike SET is_like=? WHERE user_id=? AND post_id=?`
		_, err = tx.Exec(insertSQL, isLike, userID, postID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing statement: %v", err)
		}

	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Create likes and dislikes in the database.
func LikesDelete(likeID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	insertSQL := `DELETE FROM like_dislike WHERE id=?`
	_, err = tx.Exec(insertSQL, likeID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil

}

func LikesSelectByID(likeDislikeID int64) (models.LikesDislikes, error) {
	db := SetupDatabase()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return models.LikesDislikes{}, fmt.Errorf("error starting transaction: %v", err)
	}

	// Ensure rollback in case of an error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Rethrow panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT id, user_id, post_id, comment_id, is_like, created_at
              FROM like_dislike
              WHERE id = ?`

	var likeDislike models.LikesDislikes
	err = tx.QueryRow(query, likeDislikeID).Scan(
		&likeDislike.ID,
		&likeDislike.UserID,
		&likeDislike.PostID,
		&likeDislike.CommentID,
		&likeDislike.IsLike,
		&likeDislike.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.LikesDislikes{}, fmt.Errorf("no like/dislike found with ID %d", likeDislikeID)
		}
		return models.LikesDislikes{}, fmt.Errorf("error executing query: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.LikesDislikes{}, fmt.Errorf("error committing transaction: %v", err)
	}

	return likeDislike, nil
}
