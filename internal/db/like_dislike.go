package db

import "database/sql"

func createLikesDislikesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "like_dislike" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER,
	"comment_id"	INTEGER,
	"status"	INTEGER NOT NULL,
	"createdAt"	NUMERIC DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment_id") REFERENCES "comment"("id"),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
)`
	executeSQL(db, createTableSQL)
}
