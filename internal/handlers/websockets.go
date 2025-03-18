package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
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
			"https://localhost:8080",
		}
		// Get the website link (ex: http://localhost:8080)
		origin := r.Header.Get("Origin")
		// Return true if the link is contained in the allowed links
		return slices.Contains(allowedOrigins, origin)
	},
}

// Handler that will upgrade the HTTP connection to a WebSocket connection and listen for messages
func HandleConnection(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(username, "connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(username, "disconnected")
			mu.Lock()
			delete(clients, username)
			mu.Unlock()
			break
		}

		var receivedMsg models.PrivateMessage
		err = json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			fmt.Println("Invalid JSON:", err)
			continue
		}

		if receivedMsg.Type == "private_message" {
			// sendPrivateMessage(receivedMsg)
		}
	}
}
