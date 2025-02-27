package db

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func createUsersTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "User" (
	"id"	INTEGER NOT NULL UNIQUE,
	"nickName"	TEXT NOT NULL UNIQUE,
	"gender"	TEXT NOT NULL,
	"firstName"	TEXT NOT NULL,
	"lastName"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"role"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT))`

	executeSQL(db, createTableSQL)
}

type User struct {
	ID        int
	NickName  string
	Gender    string
	FirstName string
	LastName  string
	Email     string
	Password  string // This will be hashed and not returned in most queries
	Role      string
}

// Create - Register a new user
func UserInsert(nickName, gender, firstName, lastName, email, password, role string) (int, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %v", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error hashing password: %v", err)
	}

	// Insert user
	createSQL := `INSERT INTO User (nickName, gender, firstName, lastName, email, password, role) 
                 VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := tx.Exec(createSQL, nickName, gender, firstName, lastName, email, string(hashedPassword), role)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error getting last inserted user ID: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("error committing transaction: %v", err)
	}

	return int(userID), nil
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
func UserSelectByCredentials(login string) (*User, error) {
	db := SetupDatabase()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	query := `SELECT id, nickName, gender, firstName, lastName, email, password, role 
             FROM User WHERE nickName = ? OR email = ?`

	var user User
	err = tx.QueryRow(query, login, login).Scan(
		&user.ID, &user.NickName, &user.Gender, &user.FirstName,
		&user.LastName, &user.Email, &user.Password, &user.Role,
	)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with credentials: %s", login)
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &user, nil
}

// Authenticate user
func UserAuthenticate(login, password string) (*User, error) {
	user, err := UserSelectByCredentials(login)
	if err != nil {
		return nil, err
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
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
