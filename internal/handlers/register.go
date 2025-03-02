package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // For development; add proper origin checks in production
	},
}

// Message structures
type RegisterRequest struct {
	Type      string `json:"type"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type RegisterResponse struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Storing the db value into a variable
	db := db.SetupDatabase()
	defer db.Close()

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	log.Println("Client connected")

	for {
		// Read message from client
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Parse message
		var req RegisterRequest
		if err := json.Unmarshal(message, &req); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}
		fmt.Println("Request: ", req)

		// Handle registration
		if req.Type == "register" {
			response := RegisterResponse{Type: "register_response"}

			// Check if username already exists
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM user WHERE nickName = ?", req.Username).Scan(&count)
			if err != nil {
				log.Println("Database query error:", err)
				response.Success = false
				response.Message = "Server error, please try again"
			} else if count > 0 {
				response.Success = false
				response.Message = "Username already exists"
			} else {
				// Insert new user
				_, err := db.Exec(
					"INSERT INTO user (firstName, lastName, nickName, gender, password, email, role) VALUES (?, ?, ?, ?, ?, ?, ?)",
					req.Firstname, req.Lastname, req.Username, req.Gender, req.Password, req.Email, "User") // In production, hash the password!

				if err != nil {
					log.Println("Insert error:", err)
					response.Success = false
					response.Message = "Failed to register user"
				} else {
					response.Success = true
					response.Message = "Registration successful"
				}
			}

			fmt.Println("Response: ", response)

			// Send response
			responseJSON, _ := json.Marshal(response)
			if err := ws.WriteMessage(messageType, responseJSON); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}
