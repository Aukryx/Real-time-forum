package config

import (
	"os"
	"path/filepath"
)

var (
	// Update this to point to the actual location of your forum.db file
	DB_PATH = filepath.Join("internal", "db", "forum.db")
	DB_USER = "admin"
	DB_PW   = "password"
)

// Initialize function to validate and create necessary paths
func Initialize() {
	// Ensure the database directory exists
	dbDir := filepath.Dir(DB_PATH)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0755)
	}
}
