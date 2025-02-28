package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

func createImagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "images" (
		"id"		INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER,
		file_path 	TEXT NOT NULL,
		file_size 	INTEGER NOT NULL,
		created_at	NUMERIC DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("post_id") REFERENCES "post"("id") ON DELETE CASCADE
);
`

	executeSQL(db, createTableSQL)
}

func ImageInsert(postID int, filePath string, fileSize int) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	createSQL := `INSERT INTO images (post_id, file_path, file_size) 
                  VALUES (?, ?, ?)`
	result, err := tx.Exec(createSQL, postID, filePath, fileSize)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	imageID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error getting last inserted image ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return int(imageID), nil
}

// Read - Get image by ID
func ImageSelectByID(imageID int) (*models.Image, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, post_id, file_path, file_size, created_at
              FROM images WHERE id = ?`

	var image models.Image
	var createdAtStr string

	err = tx.QueryRow(query, imageID).Scan(
		&image.ID, &image.PostID, &image.FilePath, &image.FileSize, &createdAtStr,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no image found with ID %d", imageID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Parse time string
	image.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &image, nil
}

// Read - Get images by post ID
func ImageSelectByPostID(postID int) ([]*models.Image, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, post_id, file_path, file_size, created_at
              FROM images WHERE post_id = ?`

	rows, err := tx.Query(query, postID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var images []*models.Image
	for rows.Next() {
		image := &models.Image{}
		var createdAtStr string

		if err := rows.Scan(&image.ID, &image.PostID, &image.FilePath,
			&image.FileSize, &createdAtStr); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning image: %v", err)
		}

		// Parse time string
		image.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)

		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating images: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return images, nil
}

// Update - Update image file path
func ImageUpdateFilePath(imageID int, filePath string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	updateSQL := `UPDATE images SET file_path = ? WHERE id = ?`
	_, err = tx.Exec(updateSQL, filePath, imageID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete image
func ImageDelete(imageID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM images WHERE id = ?`
	_, err = tx.Exec(deleteSQL, imageID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
