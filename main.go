package main

import (
	"log"
	"net/http"
)

func main() {
	dbConnect()
	defer dbDisconnect()
	err := dbSetup()
	if err != nil {
		log.Fatalf("Failed to access database table: %s", err)
	}

	r := http.NewServeMux()

	// TODO: add logger
	// TODO: add authentication
	// TODO: send correct status code for each request
	// TODO: improve error handling

	r.HandleFunc("/{$}", homeHandler)
	r.HandleFunc("GET /notes/{id}", getNoteHandler)
	r.HandleFunc("POST /notes", createNoteHandler)
	r.HandleFunc("DELETE /notes/{id}", deleteNoteHandler)
	r.HandleFunc("PUT /notes/{id}", updateNoteHandler)

	http.ListenAndServe("localhost:3000", r)
}
