package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"config"

	// Import the SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func SetupDatabase() *sql.DB {
	// Check if config is initialized
	config.Initialize()

	/***********************************************************************
	* Build the connection string activating authentication and encryption.
	* It'll only work using go tags like:
	* go run -tags sqlite_userauth cmd/golang-server-layout/main.go
	* You can use the tag to go test the authentication auth_test.go
	/**********************************************************************/
	connString := fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256",
		config.DB_PATH, config.DB_USER, config.DB_PW)

	// Check if the database file exists
	dbExists := false
	if _, err := os.Stat(config.DB_PATH); err == nil {
		dbExists = true
	}

	// Open or create the database file
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		log.Fatal(err)
	}

	// Set global DB variable
	DB = db

	// Only create tables if the database doesn't exist
	if !dbExists {
		fmt.Println("Creating new database tables...")
		createUsersTable(db)
		createPostsTable(db)
		createCommentsTable(db)
		createLikesDislikesTable(db)
		createImagesTable(db)
		createNotificationsTable(db)
		createPrivateMessageTable(db)
	}

	return db
}

// executeSQL prepares and executes a given SQL statement.
// It logs any errors that occur during preparation or execution.
func executeSQL(db *sql.DB, sql string) {
	// Prepare the SQL statement for execution
	statement, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err) // Log and terminate if there is an error preparing the statement
	}
	defer statement.Close()

	// Execute the prepared statement
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err) // Log and terminate if there is an error executing the statement
	}
}
