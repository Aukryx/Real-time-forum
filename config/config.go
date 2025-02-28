package config

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	DB_PATH string
	DB_USER = "admin"
	DB_PW   = "password"
)

// Initialize function to validate and create necessary paths
func Initialize() {
	// Get the absolute path to the project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")

	// Set DB_PATH to the absolute path
	DB_PATH = filepath.Join(projectRoot, "internal", "db", "forum.db")

	// Ensure the database directory exists
	dbDir := filepath.Dir(DB_PATH)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.MkdirAll(dbDir, 0755)
	}
}
