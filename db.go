package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type post struct {
	id      int
	author  string
	message string
}

var db *sql.DB

func dbConnect() error {
	connection, err := sql.Open("sqlite3", "./annynotes.db")
	if err != nil {
		log.Fatalf("Failed to access database: %s", err)
	}

	db = connection

	return nil
}

func dbDisconnect() error {
	return db.Close()
}

// Initialize database
func dbSetup() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS posts (id INTEGER NOT NULL PRIMARY KEY, author TEXT, message TEXT);
	INSERT INTO posts (author, message) VALUES ("Bilbo Baggins", "Let the adventure begin...");`)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func getPosts() ([]post, error) {
	var posts []post

	rows, err := db.Query("SELECT id, author, message FROM posts ORDER BY id DESC;")
	if err != nil {
		log.Fatalf("Failed to query database for posts: %s", err)
	}

	for rows.Next() {
		post := post{}
		rows.Scan(&post.id, &post.author, &post.message)
		posts = append(posts, post)
	}

	return posts, nil
}
