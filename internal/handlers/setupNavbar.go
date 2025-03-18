package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
)

// Response structure
type Response struct {
	Username string `json:"username"`
}

// Function to get the user's informations for the navbar display
func NavbarHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the method is correct
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieving the username with the uuid stored in the cookie
	var username string
	cookie, errCookie := r.Cookie("session_id")
	if errCookie != nil {
		fmt.Println("Error retrieving cookie in navbarhandler: ", errCookie)
		return
	}
	username = db.UserNicknameWithUUID(cookie.Value)

	// Create the response
	response := Response{
		Username: username,
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the JSON response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
