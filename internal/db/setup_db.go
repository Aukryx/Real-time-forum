package db

import (
	"config"
	"database/sql"
	"fmt"
	"log"

	// Import the SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() *sql.DB {
	/***********************************************************************
	* Build the connexion string activating authentification and crypting.
	* It'll only work  usign go tags like :
	* go run -tags sqlite_userauth cmd/golang-server-layout/main.go
	* You can use the tag to go test the authentification auth_test.go
	/**********************************************************************/
	connString := fmt.Sprintf("%s?_auth&_auth_user=%s&_auth_pass=%s&_auth_crypt=sha256", config.DB_PATH, config.DB_USER, config.DB_PW)
	// Open or create the database file
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		log.Fatal(err)
	}

	// Create all tables
	createUsersTable(db)
	createPostsTable(db)
	createCommentsTable(db)
	createLikesDislikesTable(db)
	createImagesTable(db)
	createNotificationsTable(db)
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

	// Execute the prepared statement
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err) // Log and terminate if there is an error executing the statement
	}
}
