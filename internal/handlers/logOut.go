package handlers

import (
	"db"
	"fmt"
	"net/http"
	"time"
)

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	db := db.SetupDatabase()

	// Checking the cookie values
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("Error accessing cookie: ", err)
	}

	// Unlogging the User in the database
	state := `UPDATE user SET connected = ? WHERE uuid = ?`
	_, err_db := db.Exec(state, 0, cookie.Value)
	if err_db != nil {
		fmt.Printf("Error logging out")
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expire immédiatement
		HttpOnly: true,
		Secure:   true, // Mets false si tu es en développement sans HTTPS
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
