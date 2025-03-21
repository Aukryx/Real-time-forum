// SendChatHistory retrieves message history between two users
package db

import (
	"encoding/json"
	"fmt"
	"models"

	"github.com/gorilla/websocket"
)

func SendChatHistory(user1ID, user2ID int, conn *websocket.Conn) error {
	fmt.Println("Debug: Starting SendChatHistory for users", user1ID, "and", user2ID)

	// Getting the usernames of the two users
	user1Name := UserNicknameWithID(user1ID)
	user2Name := UserNicknameWithID(user2ID)

	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// fmt.Println("Debug: Transaction error:", err)
		return fmt.Errorf("error starting transaction: %v", err)
	}
	// fmt.Println("Debug: Transaction started successfully")

	// Note: Using "user" instead of "User" since that's the table name in your schema
	query := `SELECT pm.sender_id, u_sender.nickName, pm.receiver_id, u_receiver.nickName, 
			pm.message, pm.createdAt, pm.read 
			FROM private_message pm
			JOIN user u_sender ON pm.sender_id = u_sender.id
			JOIN user u_receiver ON pm.receiver_id = u_receiver.id
			WHERE (pm.sender_id = ? AND pm.receiver_id = ?) 
			OR (pm.sender_id = ? AND pm.receiver_id = ?)
			ORDER BY pm.createdAt ASC`

	// fmt.Println("Debug: Executing query:", query)
	// fmt.Println("Debug: With params:", user1ID, user2ID, user2ID, user1ID)

	rows, err := tx.Query(query, user1ID, user2ID, user2ID, user1ID)
	if err != nil {
		// fmt.Println("Debug: Query error:", err)
		tx.Rollback()
		return fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()
	// fmt.Println("Debug: Query executed successfully")

	var messages []models.ChatHistoryMessage
	messageCount := 0

	for rows.Next() {
		var senderID, receiverID int
		var senderUsername, receiverUsername, message, createdAt string
		var read int

		err := rows.Scan(&senderID, &senderUsername, &receiverID, &receiverUsername,
			&message, &createdAt, &read)

		if err != nil {
			fmt.Println("Debug: Scan error:", err)
			tx.Rollback()
			return fmt.Errorf("error scanning message: %v", err)
		}

		messageCount++
		// fmt.Printf("Debug: Scanned message %d: sender=%s, content=%s, time=%s\n",
		// messageCount, senderUsername, message, createdAt)

		// Convert the data to the appropriate format
		chatMsg := models.ChatHistoryMessage{
			Type:      "chat_history_message",
			SenderID:  senderID,
			Sender:    senderUsername,
			Message:   message,
			Timestamp: createdAt,
			Read:      read != 0,
		}

		messages = append(messages, chatMsg)
	}

	if err = rows.Err(); err != nil {
		// fmt.Println("Debug: Row iterator error:", err)
		tx.Rollback()
		return fmt.Errorf("error iterating messages: %v", err)
	}

	fmt.Printf("Debug: Found %d messages between users\n", messageCount)

	if err = tx.Commit(); err != nil {
		// fmt.Println("Debug: Commit error:", err)
		return fmt.Errorf("error committing transaction: %v", err)
	}
	// fmt.Println("Debug: Transaction committed successfully")

	// Create the response containing the full chat history
	response := models.ChatHistory{
		Type:      "chat_history",
		User1Name: user1Name,
		User2Name: user2Name,
		Messages:  messages,
	}

	fmt.Println("response: ", response)

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		// fmt.Println("Debug: JSON marshal error:", err)
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	// fmt.Println("Debug: Response marshaled successfully")

	// Send the chat history to the requesting client
	if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
		// fmt.Println("Debug: WebSocket write error:", err)
		return fmt.Errorf("error sending chat history: %v", err)
	}
	fmt.Println("Debug: Chat history sent successfully")

	return nil
}
