package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	notes, _ := getNotes()

	data, err := json.Marshal(&notes)
	if err != nil {
		log.Fatalf("Failed to marshal notes: %s", err)
	}

	w.Write(data)
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	note, err := getNote(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to delete note with id %q: %s", id, err)
	}

	w.Write([]byte(note.Author + ": " + note.Message))
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	note := Note{}

	err := json.Unmarshal(body, &note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to unmarshal request body: %s", err)
	}

	createNote(note)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(note.Author + ": " + note.Message))
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := deleteNote(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to delete note with id %q: %s", id, err)
	}
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	newNote := Note{}

	err := json.Unmarshal(body, &newNote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to unmarshal request body: %s", err)
	}

	err = updateNote(id, newNote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to update note with id %q: %s", id, err)
	}
}
