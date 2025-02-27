package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

func createNotificationsTable(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS notification (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		sender_id INTEGER NOT NULL,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		related_id INTEGER NOT NULL,
		read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (sender_id) REFERENCES users(id)
	);`
	executeSQL(db, sql)
}

func NotificationInsert(userID int, senderID int, notificationType string, content string, relatedID int) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	createSQL := `INSERT INTO notification (user_id, sender_id, type, content, related_id) 
                  VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(createSQL, userID, senderID, notificationType, content, relatedID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	notificationID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error getting last inserted notification ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return int(notificationID), nil
}

// Read - Get notification by ID
func NotificationSelectByID(notificationID int) (*models.Notification, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, sender_id, type, content, read, related_id, createdAt
              FROM notification WHERE id = ?`

	var notification models.Notification
	var createdAtStr string
	var readInt int
	var senderID, relatedID sql.NullInt64

	err = tx.QueryRow(query, notificationID).Scan(
		&notification.ID, &notification.UserID, &senderID, &notification.Type,
		&notification.Content, &readInt, &relatedID, &createdAtStr,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no notification found with ID %d", notificationID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Handle nullable fields
	if senderID.Valid {
		notification.SenderID = int(senderID.Int64)
	}
	if relatedID.Valid {
		notification.RelatedID = int(relatedID.Int64)
	}

	// Parse time string and boolean
	notification.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	notification.Read = readInt != 0

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &notification, nil
}

// Read - Get notifications by user ID
func NotificationSelectByUserID(userID int) ([]*models.Notification, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, user_id, sender_id, type, content, read, related_id, createdAt
              FROM notification WHERE user_id = ?
              ORDER BY createdAt DESC`

	rows, err := tx.Query(query, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notification := &models.Notification{}
		var createdAtStr string
		var readInt int
		var senderID, relatedID sql.NullInt64

		if err := rows.Scan(&notification.ID, &notification.UserID, &senderID,
			&notification.Type, &notification.Content, &readInt,
			&relatedID, &createdAtStr); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning notification: %v", err)
		}

		// Handle nullable fields
		if senderID.Valid {
			notification.SenderID = int(senderID.Int64)
		}
		if relatedID.Valid {
			notification.RelatedID = int(relatedID.Int64)
		}

		// Parse time string and boolean
		notification.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		notification.Read = readInt != 0

		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating notifications: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return notifications, nil
}

// Update - Mark notification as read
func NotificationUpdateReadStatus(notificationID int, read bool) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	readInt := 0
	if read {
		readInt = 1
	}

	updateSQL := `UPDATE notification SET read = ? WHERE id = ?`
	_, err = tx.Exec(updateSQL, readInt, notificationID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete notification
func NotificationDelete(notificationID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM notification WHERE id = ?`
	_, err = tx.Exec(deleteSQL, notificationID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
