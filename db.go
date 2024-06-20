package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS posts (id INTEGER NOT NULL PRIMARY KEY, author TEXT, message TEXT);`)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func getPosts() ([]Post, error) {
	var posts []Post

	rows, err := db.Query("SELECT id, author, message FROM posts ORDER BY id ASC;")
	if err != nil {
		log.Fatalf("Failed to query database for posts: %s", err)
	}

	for rows.Next() {
		post := Post{}
		rows.Scan(&post.ID, &post.Author, &post.Message)
		posts = append(posts, post)
	}

	return posts, nil
}

func getPost(id string) (Post, error) {
	post := Post{}

	err := db.QueryRow(`SELECT author, message FROM posts WHERE id = ?`, id).Scan(&post.Author, &post.Message)
	if err != nil {
		log.Panic(err)
	}

	return post, nil
}

func createPost(p Post) error {
	_, err := db.Exec(`INSERT INTO posts (author, message) VALUES (?, ?);`, p.Author, p.Message)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func deletePost(id string) error {
	_, err := db.Exec(`DELETE FROM posts WHERE id = ?;`, id)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func updatePost(id string, p Post) error {
	post := Post{}
	err := db.QueryRow(`SELECT author, message FROM posts WHERE id = ?`, id).Scan(&post.Author, &post.Message)
	if err != nil {
		log.Panic(err)
	}

	if p.Author != "" {
		post.Author = p.Author
	}

	if p.Message != "" {
		post.Message = p.Message
	}

	_, err = db.Exec(`UPDATE posts SET author = ?, message = ? WHERE id = ?;`, post.Author, post.Message, id)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
