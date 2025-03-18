package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"models"
	"time"

	"github.com/gorilla/websocket"
)

var clients = models.GetClientMap()
var mu = models.GetMux()

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

func sendPrivateMessage(msg models.PrivateMessage) {
	// Ensure both sender and receiver are set
	if msg.Sender == "" || msg.Receiver == "" {
		fmt.Println("Error: Sender or receiver not specified")
		return
	}

	// Check if the receiver exists (is connected)
	mu.Lock()
	receiverConn, receiverExists := clients[msg.Receiver]
	senderConn, senderExists := clients[msg.Sender]
	mu.Unlock()

	// Create a response message
	response := models.PrivateMessage{
		Type:     "private_message",
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
		Message:  msg.Message,
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Send message to receiver if they're connected
	if receiverExists {
		err := receiverConn.WriteMessage(websocket.TextMessage, jsonResponse)
		if err != nil {
			fmt.Println("Error sending message to receiver:", err)
		}
	} else {
		// Store the message in the database for offline delivery
		// db.StoreOfflineMessage(msg.Sender, msg.Receiver, msg.Message)

		// Notify sender that receiver is offline
		if senderExists {
			notifyMsg := models.PrivateMessage{
				Type:     "system_notification",
				Sender:   "system",
				Receiver: msg.Sender,
				Message:  msg.Receiver + " is currently offline. Message will be delivered when they connect.",
			}

			notifyJson, _ := json.Marshal(notifyMsg)
			senderConn.WriteMessage(websocket.TextMessage, notifyJson)
		}
	}

	// Also send a copy/confirmation to the sender
	if senderExists {
		confirmMsg := models.PrivateMessage{
			Type:     "message_sent",
			Sender:   msg.Sender,
			Receiver: msg.Receiver,
			Message:  msg.Message,
		}

		confirmJson, _ := json.Marshal(confirmMsg)
		senderConn.WriteMessage(websocket.TextMessage, confirmJson)
	}

	// Log the message
	fmt.Printf("Private message from %s to %s: %s\n", msg.Sender, msg.Receiver, msg.Message)
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
		&message.Type, &message.Sender, &message.Receiver, &message.Message,
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

		if err := rows.Scan(&message.Type, &message.Sender, &message.Receiver,
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
