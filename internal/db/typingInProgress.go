package db

import (
	"encoding/json"
	"fmt"
	"models"

	"github.com/gorilla/websocket"
)

// Function to notify user when someone is typing
func TypingInProgress(msg models.PrivateMessage) {
	// Ensure both sender and receiver are set
	if msg.Sender == "" || msg.Receiver == "" {
		fmt.Println("Error: Sender or receiver not specified")
		return
	}

	// Check if the receiver exists (is connected)
	mu.Lock()
	receiverConn, receiverExists := clients[msg.Receiver]
	mu.Unlock()

	// Create a response message
	response := models.PrivateMessage{
		Type:     "typing",
		Sender:   msg.Sender,
		Receiver: msg.Receiver,
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
	}

	// Log the message (not recommanded xd)
	// fmt.Printf("Private message from %s to %s: %s\n", msg.Sender, msg.Receiver)
}
