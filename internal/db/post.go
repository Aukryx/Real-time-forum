package db

import "database/sql"

func createPostsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "post" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"updatedAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"image"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}
