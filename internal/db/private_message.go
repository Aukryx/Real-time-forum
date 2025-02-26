package db

import "database/sql"

func createPrivateMessageTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "private_message" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id1"	INTEGER NOT NULL,
	"user_id2"	INTEGER NOT NULL,
	"body1"	TEXT,
	"body2"	TEXT,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id1") REFERENCES "User"("id"),
	FOREIGN KEY("user_id2") REFERENCES "User"("id")
)`

	executeSQL(db, createTableSQL)
}
