BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "User" (
	"id"	INTEGER NOT NULL UNIQUE,
	"nickName"	TEXT NOT NULL UNIQUE,
	"gender"	TEXT NOT NULL,
	"firstName"	TEXT NOT NULL,
	"lastName"	TEXT NOT NULL,
	"email"	TEXT NOT NULL UNIQUE,
	"password"	TEXT NOT NULL,
	"role"	TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "comment" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC,
	"updatedAt"	NUMERIC,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
);
CREATE TABLE IF NOT EXISTS "like_dislike" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER,
	"comment_id"	INTEGER,
	"status"	INTEGER NOT NULL,
	"createdAt"	NUMERIC,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment_id") REFERENCES "comment"("id"),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
);
CREATE TABLE IF NOT EXISTS "notification" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"post_id"	INTEGER,
	"comment_id"	INTEGER,
	"createdAt"	NUMERIC,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("comment_id") REFERENCES "comment"("id"),
	FOREIGN KEY("post_id") REFERENCES "post"("id"),
	FOREIGN KEY("user_id") REFERENCES "User"("id")
);
CREATE TABLE IF NOT EXISTS "post" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"body"	TEXT NOT NULL,
	"createdAt"	NUMERIC NOT NULL,
	"updatedAt"	NUMERIC NOT NULL,
	"image"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("id") REFERENCES "",
	FOREIGN KEY("user_id") REFERENCES "User"("id")
);
CREATE TABLE IF NOT EXISTS "private_message" (
	"id"	INTEGER NOT NULL UNIQUE,
	"user_id1"	INTEGER NOT NULL,
	"user_id2"	INTEGER NOT NULL,
	"body1"	TEXT,
	"body2"	TEXT,
	"createdAt"	NUMERIC NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("user_id1") REFERENCES "User"("id"),
	FOREIGN KEY("user_id2") REFERENCES "User"("id")
);
COMMIT;
