package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"os"
	"slices"

	"github.com/gorilla/websocket"
)

var clients = models.GetClientMap()
var mu = models.GetMux()

// Upgrader to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow WebSocket connections from http://localhost:8080
		allowedOrigins := []string{
			"http://localhost:8080",
			"http://localhost:8081",
			"http://localhost:8082",
			"http://localhost:8040",
			"http://localhost:8050",
			"http://localhost:8060",
			"http://localhost:8070",
			"https://real-time-forum.onrender.com",
		}
		// Get the website link (ex: http://localhost:8080)
		origin := r.Header.Get("Origin")
		// Return true if the link is contained in the allowed links
		return slices.Contains(allowedOrigins, origin)
	},
}

// Handler that will upgrade the HTTP connection to a WebSocket connection and listen for messages
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	// Check if running on Render
	if os.Getenv("PORT") != "" {
		// Return a friendly message instead of attempting WebSocket connection
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "disabled",
			"message": "WebSockets are disabled in the demo version. Please check the GitHub repository for the full application.",
		})
		return
	}

	// Regular WebSocket handling code for local development
	// Retrieving the cookie for the uuid
	cookie, errCookie := r.Cookie("session_id")
	if errCookie != nil {
		fmt.Println("Error getting cookie:", errCookie)
		return
	}
	// Getting the username with the UUID stored in the cookie
	username := db.UserNicknameWithUUID(cookie.Value)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[username] = conn
	mu.Unlock()
	SendUserListToAll()

	fmt.Println(username, "connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(username, "disconnected")
			mu.Lock()
			delete(clients, username)
			mu.Unlock()
			SendUserListToAll()
			break
		}

		var receivedMsg models.PrivateMessage
		err = json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		// Get the sender and receiver IDs
		sender := db.UserIDWithNickname(receivedMsg.Sender)
		receiver := db.UserIDWithNickname(receivedMsg.Receiver)

		// Check the type of message
		if receivedMsg.Type == "private_message" {
			fmt.Println("Received private message from", receivedMsg.Sender, "to", receivedMsg.Receiver)

			// Send the message to the receiver, client side
			db.SendPrivateMessage(receivedMsg)

			// Insert the message into the database
			db.PrivateMessageInsert(sender, receiver, receivedMsg.Message)
		} else if receivedMsg.Type == "chat_history_request" {
			fmt.Println("Received chat history request between", receivedMsg.Sender, "to", receivedMsg.Receiver)
			db.SendChatHistory(sender, receiver, conn)
		}
	}
}

// SendUserListToAll sends the current list of connected users to all connected clients
func SendUserListToAll() {
	mu.Lock()
	// Create a slice of all connected usernames
	userList := make([]string, 0, len(clients))
	for username := range clients {
		userList = append(userList, username)
	}

	// Prepare the message once
	response := models.PrivateMessage{
		Type:     "user_list",
		Sender:   "server",
		Message:  "",
		UserList: userList,
	}

	// Convert to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling user list:", err)
		mu.Unlock()
		return
	}

	// Send to all connected clients
	for username, conn := range clients {
		response.Receiver = username // Set the receiver for logging purposes
		err := conn.WriteMessage(websocket.TextMessage, jsonResponse)
		if err != nil {
			fmt.Println("Error sending user list to", username, ":", err)
		}
	}
	mu.Unlock()
}
