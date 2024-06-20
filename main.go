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
	r.HandleFunc("GET /posts/{id}", getPostHandler)
	r.HandleFunc("POST /posts", createPostHandler)
	r.HandleFunc("DELETE /posts/{id}", deletePostHandler)
	r.HandleFunc("PUT /posts/{id}", updatePostHandler)

	http.ListenAndServe("localhost:3000", r)
}
