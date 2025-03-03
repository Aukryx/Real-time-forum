package handlers

import (
	"db"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Name     string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ensuring the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Getting the JSON form data to test
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Inserting the user into the database
	_, errorDB := db.UserAuthenticate(req.Name, req.Password)

	// Checking if the insert failed
	if errorDB != "nil" {
		response := RegisterResponse{Success: false, Message: errorDB}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// If the insert didn't fail, notify the js of the success
	json.NewEncoder(w).Encode(RegisterResponse{Success: true, Message: "Login successful"})
}
