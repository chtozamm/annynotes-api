package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/chtozamm/annynotes-go/internal/database"
	"github.com/chtozamm/annynotes-go/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

// The core of the application
type application struct {
	DB *database.Queries
}

func main() {
	// Load envinronmental variables from .env
	utils.ParseEnv(".env")
	port := os.Getenv("PORT")
	if port == "" {
		log.Print("No PORT variable was found in .env, default value is set")
		port = "3000"
	} else {
		log.Printf("Found PORT variable in .env: %s", port)
	}

	// Set a database file location depending on the OS
	var localDataPath string
	if runtime.GOOS == "windows" {
		localDataPath = os.Getenv("LOCALAPPDATA")
	} else {
		localDataPath = os.Getenv("XDG_DATA_HOME")
	}
	// Open database file or create a new one if doesn't exist
	connection, err := sql.Open("sqlite3", path.Join(localDataPath, "annynotes", "annynotes.db"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err)
	}
	// Setup database tables
	err = setupDB(connection)
	if err != nil {
		log.Fatalf("Failed to setup the database: %s", err)
	}

	app := application{
		DB: database.New(connection),
	}

	r := http.NewServeMux()

	r.HandleFunc("/{$}", app.homeHandler)
	r.HandleFunc("GET /notes/{id}", app.getNoteHandler)
	r.HandleFunc("POST /notes", app.withAuth(app.createNoteHandler))
	r.HandleFunc("PUT /notes/{id}", app.withAuth(app.updateNoteHandler))
	r.HandleFunc("DELETE /notes/{id}", app.withAuth(app.deleteNoteHandler))
	r.HandleFunc("GET /{author}", app.getNotesFromAuthorHandler)
	r.HandleFunc("POST /users", app.createUserHandler)
	r.HandleFunc("POST /users/auth", app.authenticateUserHandler)

	log.Print("Server is listening on localhost:", port)
	http.ListenAndServe("localhost:"+port, r)
}
