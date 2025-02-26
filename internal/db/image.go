package db

import (
	"database/sql"
)

func createImagesTable(db *sql.DB) {
	createTableSQL := `CREATE TABLE IF NOT EXISTS "images" (
		"id"		INTEGER NOT NULL UNIQUE,
		"post_id"	INTEGER,
		file_path 	TEXT NOT NULL,
		file_size 	INTEGER NOT NULL,
		created_at	NUMERIC DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY("id" AUTOINCREMENT),
		FOREIGN KEY("post_id") REFERENCES "post"("id") ON DELETE CASCADE
);
`

	executeSQL(db, createTableSQL)
}
