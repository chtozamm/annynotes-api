package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS notes (
	id INTEGER NOT NULL PRIMARY KEY,
  author TEXT,
  message TEXT,
  created_at TEXT DEFAULT CURRENT_TIMESTAMP
);`)
	if err != nil {
		log.Fatalf("Failed to setup database: %s", err)
	}

	return nil
}

func getNotes() ([]Note, error) {
	var notes []Note

	rows, err := db.Query("SELECT id, author, message FROM notes ORDER BY id ASC;")
	if err != nil {
		log.Fatalf("Failed to query database for notes: %s", err)
	}

	for rows.Next() {
		note := Note{}
		rows.Scan(&note.ID, &note.Author, &note.Message)
		notes = append(notes, note)
	}

	return notes, nil
}

func getNotesDesc() ([]Note, error) {
	var notes []Note

	rows, err := db.Query("SELECT * FROM notes ORDER BY id DESC;")
	if err != nil {
		log.Fatalf("Failed to query database for notes: %s", err)
	}

	for rows.Next() {
		note := Note{}
		rows.Scan(&note.ID, &note.Author, &note.Message, &note.CreatedAt)
		notes = append(notes, note)
	}

	return notes, nil
}

func getNote(id string) (Note, error) {
	note := Note{}

	err := db.QueryRow(`SELECT * FROM notes WHERE id = ?;`, id).Scan(&note.ID, &note.Author, &note.Message, &note.CreatedAt)
	if err != nil {
		log.Panic(err)
	}

	return note, nil
}

func getAmountOfNotes() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM notes;").Scan(&count)
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func createNote(n Note) (Note, error) {
	note := Note{}
	err := db.QueryRow(`INSERT INTO notes (author, message) VALUES (?, ?) RETURNING id, author, message;`, n.Author, n.Message).Scan(&note.ID, &note.Author, &note.Message)
	if err != nil {
		log.Panic(err)
	}

	return note, nil
}

func deleteNote(id string) error {
	_, err := db.Exec(`DELETE FROM notes WHERE id = ?;`, id)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func updateNote(id string, n Note) error {
	note := Note{}
	err := db.QueryRow(`SELECT author, message FROM notes WHERE id = ?`, id).Scan(&note.Author, &note.Message)
	if err != nil {
		log.Panic(err)
	}

	if n.Author != "" {
		note.Author = n.Author
	}

	if n.Message != "" {
		note.Message = n.Message
	}

	_, err = db.Exec(`UPDATE notes SET author = ?, message = ? WHERE id = ?;`, note.Author, note.Message, id)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
