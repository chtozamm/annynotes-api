package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getNotes() ([]Note, error) {
	resp, err := http.Get("http://localhost:3000")
	if err != nil {
		log.Fatalf("Failed to get http response: %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s", err)
	}
	defer resp.Body.Close()

	notes := []Note{}

	err = json.Unmarshal(body, &notes)
	if err != nil {
		log.Fatalf("Failed to unmarshal notes: %s", err)
	}

	return notes, nil
}
