package handlers

import (
	"db"
	"encoding/json"
	"middlewares"
	"models"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Ensuring the method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Getting the JSON form data to test
	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	uuid := middlewares.GenerateSessionID()

	// Inserting the user into the database
	userID, errorMsg := db.UserInsert(uuid, req.Username, req.Gender, req.Firstname, req.Lastname, req.Email, req.Password, "User", 1)

	// Checking if the insert failed
	if userID == 0 {
		response := models.RegisterResponse{Success: false, Message: errorMsg}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Maintenant que l'utilisateur est enregistré, créer une session
	middlewares.CreateSession(w, userID, req.Username, "User", uuid)

	// If the insert didn't fail, notify the js of the success
	json.NewEncoder(w).Encode(models.RegisterResponse{Success: true, Message: "Registration successful"})
}
