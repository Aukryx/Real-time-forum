package db

import "database/sql"

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
	PRIMARY KEY("id" AUTOINCREMENT)`

	// Call executeSQL to run the SQL statement and create the table
	executeSQL(db, createTableSQL)
}
