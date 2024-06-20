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

	r.HandleFunc("/", homeHandler)

	http.ListenAndServe("localhost:3000", r)
}
