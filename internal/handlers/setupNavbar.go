package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
)

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
		// Return a response with empty username or some default state
		response := models.Response{
			Username: "", // or "Guest" or whatever makes sense for your application
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return // Important: return here to avoid using cookie.Value when cookie is nil
	}

	// Only try to get the username if we have a valid cookie
	username = db.UserNicknameWithUUID(cookie.Value)

	// Create the response
	response := models.Response{
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
