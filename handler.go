package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note

	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "desc" {
		notes, _ = getNotesDesc()
	} else {
		notes, _ = getNotes()
	}

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

func getCountHandler(w http.ResponseWriter, r *http.Request) {
	count, err := getAmountOfNotes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Failed to fetch amount of notes: %s", err)
	}

	countStr := strconv.Itoa(count)

	w.Write([]byte(countStr))
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

	n, err := createNote(note)
	if err != nil {
		log.Fatalf("Failed to create a note: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(note.Author + ": " + note.Message))
	data, err := json.Marshal(&n)
	if err != nil {
		log.Fatalf("Failed to marshal notes: %s", err)
	}

	w.Write(data)
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
