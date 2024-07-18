package main

import (
	"database/sql"
	"fmt"
)

func dbConnect(path string) (*sql.DB, error) {
	// Open database file or create a new one if doesn't exist
	connection, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %s", err)
	}
	// Setup database tables
	err = setupDB(connection)
	if err != nil {
		return nil, fmt.Errorf("failed to setup the database: %s", err)
	}
	return connection, nil
}

func setupDB(conn *sql.DB) error {
	_, err := conn.Exec(`
CREATE TABLE IF NOT EXISTS notes (
  id TEXT NOT NULL PRIMARY KEY,
  author TEXT NOT NULL,
  message TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  user_id TEXT NOT NULL,
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_note_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
  UPDATE notes
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL CHECK(
    length(name) >= 2 AND
    length(name) <= 20
  ),
  username TEXT NOT NULL UNIQUE CHECK(
    length(username) >= 2 AND
    length(username) <= 20
  ),
  password TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_user_timestamp
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
  UPDATE users
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;
`)
	if err != nil {
		return err
	}
	return nil
}
