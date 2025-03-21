package handlers

import (
	"db"
	"encoding/json"
	"fmt"
	"middlewares"
	"models"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ensuring the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Getting the JSON form data to test
	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Authenticate the user using the database
	user, errorDB := db.UserAuthenticate(req.Name, req.Password)

	// Checking if the authentication failed
	if errorDB != "nil" {
		response := models.RegisterResponse{Success: false, Message: errorDB}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create a session for the authenticated user
	middlewares.CreateSession(w, user.ID, user.NickName, user.Role, user.UUID)

	db := db.SetupDatabase()
	// Unlogging the User in the database
	state := `UPDATE user SET connected = ? WHERE uuid = ?`
	_, err_db := db.Exec(state, 1, user.UUID)
	if err_db != nil {
		fmt.Printf("Error logging in")
	}

	// If authentication succeeded, notify the client of the success
	json.NewEncoder(w).Encode(models.RegisterResponse{Success: true, Message: "Login successful"})
}
