package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func createUsersTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "user" (
	"id"	INTEGER NOT NULL UNIQUE,
	"uuid"  TEXT NOT NULL UNIQUE,
	"nickName"	TEXT NOT NULL UNIQUE,
	"gender"	TEXT NOT NULL,
	"firstName"	TEXT NOT NULL,
	"lastName"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"role"	TEXT NOT NULL,
	"connected" INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY("id" AUTOINCREMENT))`

	executeSQL(db, createTableSQL)
}

type User struct {
	ID        int
	UUID      string
	NickName  string
	Gender    string
	FirstName string
	LastName  string
	Email     string
	Password  string // This will be hashed and not returned in most queries
	Role      string
}

// Create - Register a new user
func UserInsert(uuid, nickName, gender, firstName, lastName, email, password, role string, connected int) (int, string) {
	db := SetupDatabase()
	defer db.Close()

	// Check if username already exists
	var existingUserID int
	err := db.QueryRow("SELECT id FROM User WHERE nickName = ?", nickName).Scan(&existingUserID)
	if err == nil {
		return 0, "Username already taken"
	} else if err != sql.ErrNoRows {
		return 0, "Error checking username"
	}

	// Check if email already exists
	err = db.QueryRow("SELECT id FROM User WHERE email = ?", email).Scan(&existingUserID)
	if err == nil {
		return 0, "Email already registered"
	} else if err != sql.ErrNoRows {
		return 0, "Error checking email"
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, "Error starting transaction"
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return 0, "Error hashing password"
	}

	// Insert user
	createSQL := `INSERT INTO User (uuid, nickName, gender, firstName, lastName, email, password, role, connected) 
                 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(createSQL, uuid, nickName, gender, firstName, lastName, email, string(hashedPassword), role, connected)
	if err != nil {
		tx.Rollback()
		return 0, "Error executing query"
	}

	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, "Error getting last inserted user ID"
	}

	if err = tx.Commit(); err != nil {
		return 0, "Error committing transaction"
	}

	return int(userID), "User registered successfully"
}

// Read - Get user by ID
func UserSelectByID(userID int) (*User, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, nickName, gender, firstName, lastName, email, role 
             FROM User WHERE id = ?`

	var user User
	err = tx.QueryRow(query, userID).Scan(
		&user.ID, &user.NickName, &user.Gender, &user.FirstName,
		&user.LastName, &user.Email, &user.Role,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with ID %d", userID)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &user, nil
}

// Read - Get user by nickname or email (for login)
func UserSelectByCredentials(login string) (*User, string) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, "Error starting transaction"
	}

	query := `SELECT id, uuid, nickName, gender, firstName, lastName, email, password, role 
             FROM User WHERE nickName = ? OR email = ?`

	var user User
	err = tx.QueryRow(query, login, login).Scan(
		&user.ID, &user.UUID, &user.NickName, &user.Gender, &user.FirstName,
		&user.LastName, &user.Email, &user.Password, &user.Role,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, "No user found with those credentials"
		}
		return nil, "Error executing query"
	}

	if err = tx.Commit(); err != nil {
		return nil, "Error committing transaction"
	}

	return &user, "nil"
}

// Authenticate user
func UserAuthenticate(login, password string) (*User, string) {
	user, err := UserSelectByCredentials(login)
	if err != "nil" {
		return nil, err
	}

	// Compare password with hash
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errPassword != nil {
		return nil, "Invalid password"
	}

	// Clear password before returning
	user.Password = ""
	return user, "nil"
}

// Update - Update user information
func UserUpdate(userID int, nickName, gender, firstName, lastName, email, role string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	updateSQL := `UPDATE User SET nickName=?, gender=?, firstName=?, lastName=?, email=?, role=? 
                 WHERE id=?`
	_, err = tx.Exec(updateSQL, nickName, gender, firstName, lastName, email, role, userID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Update - Change password
func UserUpdatePassword(userID int, newPassword string) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error hashing password: %v", err)
	}

	updateSQL := `UPDATE User SET password=? WHERE id=?`
	_, err = tx.Exec(updateSQL, string(hashedPassword), userID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// Delete - Delete user
func UserDelete(userID int) error {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	deleteSQL := `DELETE FROM User WHERE id=?`
	_, err = tx.Exec(deleteSQL, userID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing statement: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// List all users
func UserSelectAll() ([]User, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, nickName, gender, firstName, lastName, email, role FROM User`

	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.NickName, &user.Gender, &user.FirstName,
			&user.LastName, &user.Email, &user.Role); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error iterating users: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return users, nil
}
