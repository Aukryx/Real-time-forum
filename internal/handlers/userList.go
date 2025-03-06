package handlers

import (
	"db"
	"encoding/json"
	"net/http"
)

func GetConnectedAndDisconnectedUsers(w http.ResponseWriter, r *http.Request) {
	db := db.SetupDatabase()
	defer db.Close()

	// Query for connected users (connected = 1)
	connectedUsersQuery := `SELECT id, nickName, gender, firstName, lastName, email, role 
	                      FROM "user" WHERE connected = 1`

	connectedRows, err := db.Query(connectedUsersQuery)
	if err != nil {
		http.Error(w, "Failed to query connected users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer connectedRows.Close()

	// Query for disconnected users (connected = 0)
	disconnectedUsersQuery := `SELECT id, nickName, gender, firstName, lastName, email, role 
	                         FROM "user" WHERE connected = 0`

	disconnectedRows, err := db.Query(disconnectedUsersQuery)
	if err != nil {
		http.Error(w, "Failed to query disconnected users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer disconnectedRows.Close()

	// Struct to represent a user (excluding password for security)
	type User struct {
		ID        int    `json:"id"`
		NickName  string `json:"nickName"`
		Gender    string `json:"gender"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Connected int    `json:"connected"`
	}

	// Response structure
	type Response struct {
		ConnectedUsers    []User `json:"connectedUsers"`
		DisconnectedUsers []User `json:"disconnectedUsers"`
	}

	// Parse connected users
	connectedUsers := []User{}
	for connectedRows.Next() {
		var user User
		user.Connected = 1
		err := connectedRows.Scan(&user.ID, &user.NickName, &user.Gender, &user.FirstName,
			&user.LastName, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Error scanning connected users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		connectedUsers = append(connectedUsers, user)
	}

	// Parse disconnected users
	disconnectedUsers := []User{}
	for disconnectedRows.Next() {
		var user User
		user.Connected = 0
		err := disconnectedRows.Scan(&user.ID, &user.NickName, &user.Gender, &user.FirstName,
			&user.LastName, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "Error scanning disconnected users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		disconnectedUsers = append(disconnectedUsers, user)
	}

	// Create the response
	response := Response{
		ConnectedUsers:    connectedUsers,
		DisconnectedUsers: disconnectedUsers,
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the JSON response
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
