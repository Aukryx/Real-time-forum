package db

import (
	"database/sql"
	"fmt"
	"models"
	"time"
)

// Note: Your current schema is questionable - having body1 and body2 fields
// suggests you're trying to store both sides of a conversation in one row,
// which would be difficult to work with. Let me provide a better approach.

func createPrivateMessageTable(db *sql.DB) {
	// A better schema for chat messages
	createTableSQL := `CREATE TABLE IF NOT EXISTS "private_message" (
	"id"	INTEGER NOT NULL UNIQUE,
	"sender_id"	INTEGER NOT NULL,
	"receiver_id"	INTEGER NOT NULL,
	"message"	TEXT NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"read"	INTEGER DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("sender_id") REFERENCES "User"("id"),
	FOREIGN KEY("receiver_id") REFERENCES "User"("id")
)`

	executeSQL(db, createTableSQL)
}

func PrivateMessageInsert(senderID, receiverID int, message string) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	createSQL := `INSERT INTO private_message (sender_id, receiver_id, message) 
                VALUES (?, ?, ?)`
	result, err := tx.Exec(createSQL, senderID, receiverID, message)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error getting last inserted message ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return int(messageID), nil
}

// Read - Get a message by ID
func PrivateMessageSelectByID(messageID int) (*models.PrivateMessage, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, sender_id, receiver_id, message, createdAt, read
              FROM private_message WHERE id = ?`

	var message models.PrivateMessage
	var createdAtStr string
	var readInt int

	err = tx.QueryRow(query, messageID).Scan(
		&message.ID, &message.SenderID, &message.ReceiverID, &message.Message,
		&createdAtStr, &readInt,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no message found with ID %d", messageID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	// Parse time string and boolean
	message.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	message.Read = readInt != 0

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &message, nil
}

// Read - Get all messages for a user (both sent and received)
func PrivateMessageSelectByUserID(userID int) ([]*models.PrivateMessage, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, sender_id, receiver_id, message, createdAt, read
              FROM private_message 
              WHERE sender_id = ? OR receiver_id = ?
              ORDER BY createdAt DESC`

	rows, err := tx.Query(query, userID, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var messages []*models.PrivateMessage
	for rows.Next() {
		message := &models.PrivateMessage{}
		var createdAtStr string
		var readInt int

		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID,
			&message.Message, &createdAtStr, &readInt); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning message: %v", err)
		}

		// Parse time string and boolean
		message.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		message.Read = readInt != 0

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating messages: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return messages, nil
}

// Update - Mark message as read
func PrivateMessageUpdateReadStatus(messageID int, read bool) error {
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

	updateSQL := `UPDATE private_message SET read = ? WHERE id = ?`
	_, err = tx.Exec(updateSQL, readInt, messageID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete message
func PrivateMessageDelete(messageID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM private_message WHERE id = ?`
	_, err = tx.Exec(deleteSQL, messageID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
