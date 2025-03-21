package db

// import (
// 	"fmt"
// 	"models"
// )

// // GetChatHistoryBetweenUsers returns paginated chat history between two users
// func GetChatHistoryBetweenUsers(user1ID, user2ID, page, pageSize int) ([]models.PrivateMessage, error) {
// 	db := SetupDatabase()
// 	defer db.Close()

// 	fmt.Printf("Fetching chat history for users %d and %d (page: %d, pageSize: %d)\n", user1ID, user2ID, page, pageSize)

// 	// Calculate offset for pagination
// 	offset := (page - 1) * pageSize

// 	// Updated query to properly join with user table and handle sender/receiver names
// 	query := `
// 		SELECT
// 			pm.id,
// 			sender.nickName as sender_name,
// 			receiver.nickName as receiver_name,
// 			pm.message,
// 			pm.createdAt
// 		FROM private_message pm
// 		JOIN User sender ON pm.sender_id = sender.id
// 		JOIN User receiver ON pm.receiver_id = receiver.id
// 		WHERE (pm.sender_id = ? AND pm.receiver_id = ?)
// 		   OR (pm.sender_id = ? AND pm.receiver_id = ?)
// 		ORDER BY pm.createdAt DESC
// 		LIMIT ? OFFSET ?
// 	`

// 	rows, err := db.Query(query, user1ID, user2ID, user2ID, user1ID, pageSize, offset)
// 	if err != nil {
// 		fmt.Printf("Error in query: %v\n", err)
// 		return nil, fmt.Errorf("error executing query: %v", err)
// 	}
// 	defer rows.Close()

// 	var messages []models.PrivateMessage
// 	messageCount := 0
// 	for rows.Next() {
// 		messageCount++
// 		var msg models.PrivateMessage
// 		var id int
// 		err := rows.Scan(&id, &msg.Sender, &msg.Receiver, &msg.Message, &msg.CreatedAt)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning message: %v", err)
// 		}

// 		msg.Type = "private_message"
// 		msg.ID = id
// 		messages = append(messages, msg)
// 	}

// 	fmt.Printf("Found %d messages in history\n", messageCount)
// 	return messages, nil
// }
