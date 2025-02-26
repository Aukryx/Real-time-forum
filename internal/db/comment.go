package db

import "database/sql"

func createCommentsTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "comment" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	"updatedAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}
