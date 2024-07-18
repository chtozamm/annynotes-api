package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/chtozamm/annynotes-go/internal/database"
	"github.com/chtozamm/annynotes-go/internal/utils"
	_ "github.com/mattn/go-sqlite3"
)

// The core of the application
type application struct {
	DB  *database.Queries
	srv *http.Server
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

	// Database
	connection, err := dbConnect(path.Join(localDataPath, "annynotes", "annynotes.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	r := http.NewServeMux()

	app := application{
		DB: database.New(connection),
		srv: &http.Server{
			Addr:    port,
			Handler: r,
		},
	}

	// Router
	r.HandleFunc("GET /notes", app.getNotesHandler)
	r.HandleFunc("POST /notes", app.withAuth(app.createNoteHandler))
	r.HandleFunc("GET /note/{id}", app.getNoteHandler)
	r.HandleFunc("PATCH /note/{id}", app.withAuth(app.updateNoteHandler))
	r.HandleFunc("DELETE /note/{id}", app.withAuth(app.deleteNoteHandler))
	r.HandleFunc("GET /notes/{author}", app.getNotesFromAuthorHandler)
	r.HandleFunc("POST /users", app.createUserHandler)
	r.HandleFunc("POST /users/auth", app.authenticateUserHandler)

	// Gracefully shut down by handling existing requests in the given time
	go func() {
		log.Print("Server is listening on localhost:", port)
		http.ListenAndServe(":"+port, r)
	}()
	// Create a channel to listen for shutdown signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// Block until a signal is received
	<-quit
	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Attempt to gracefully shutdown the server
	if err := app.srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err)
	}
	log.Println("Server closed")
}
