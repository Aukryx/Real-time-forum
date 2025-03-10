package handlers

import (
	"encoding/json"
	"net/http"
)

// Example in Go (adjust according to your actual server implementation)
func CheckSession(w http.ResponseWriter, r *http.Request) {
	// Get the session cookie
	_, err := r.Cookie("session_id")

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// No valid session
		json.NewEncoder(w).Encode(map[string]bool{"loggedIn": false})
		return
	}

	// Valid session
	json.NewEncoder(w).Encode(map[string]bool{"loggedIn": true})
}
