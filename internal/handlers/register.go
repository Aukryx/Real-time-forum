package handlers

import (
	"db"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Ensuring the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Getting the JSON form data to test
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Inserting the user into the database
	count, errorDB := db.UserInsert(req.Username, req.Gender, req.Firstname, req.Lastname, req.Email, req.Password, "User")

	// Checking if the insert failed
	if count == 0 {
		response := RegisterResponse{Success: false, Message: errorDB}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// If the insert didn't fail, notify the js of the success
	json.NewEncoder(w).Encode(RegisterResponse{Success: true, Message: "Registration successful"})
}
