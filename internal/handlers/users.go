package handlers

import (
	"db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HandleUserSelectAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := db.UserSelectAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add some logging
	log.Printf("Fetched %d users", len(users))

	w.Header().Set("Content-Type", "application/json")

	// Use json.NewEncoder for more robust encoding
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding users: %v", err)
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
}

func HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	user, err := db.UserSelectByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
