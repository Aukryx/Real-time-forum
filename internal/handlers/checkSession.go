package handlers

import (
	"database/sql"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
)

// Function to check the session with the cookie and database request
func CheckSession(w http.ResponseWriter, r *http.Request) {
	// Get the db
	db := db.SetupDatabase()

	// Get the session cookie
	cookie, err := r.Cookie("session_id")

	w.Header().Set("Content-Type", "application/json")

	// If the cookie doesn't exist, return false to the javascript
	if err != nil {
		// No valid session
		json.NewEncoder(w).Encode(map[string]bool{"loggedIn": false})
		return
	}

	// Checking if the uuid stored in the cookie is valid
	var nickname string
	state := `SELECT nickName FROM user WHERE uuid = ?`
	errQuery := db.QueryRow(state, cookie.Value).Scan(&nickname)
	fmt.Println("error: ", errQuery)

	// If it's not contained, log out the user
	if errQuery == sql.ErrNoRows {
		LogOutHandler(w, r)
	} else if errQuery != nil {
		fmt.Println("Error checking cookie session: ", errQuery)
	}

	// Valid session
	json.NewEncoder(w).Encode(map[string]bool{"loggedIn": true})
}
