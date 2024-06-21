package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func main() {
	err := parseTemplates()
	if err != nil {
		log.Fatalf("Failed to parse templates: %s", err)
	}

	r := http.NewServeMux()
	r.HandleFunc("GET /{$}", homeHandler)
	r.HandleFunc("GET /notes/show_form", showAddFormHandler)
	r.HandleFunc("GET /notes/close_form", closeAddFormHandler)
	r.HandleFunc("POST /notes", addNoteHandler)

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	notes := []Note{}
	resp, _ := http.Get("http://localhost:3000/?sort=desc")
	err := json.NewDecoder(resp.Body).Decode(&notes)
	if err != nil {
		log.Fatalf("Failed to decode response: %s", err)
	}

	type Data struct {
		Note
		Idx int
	}
	data := []Data{}

	idx := len(notes)
	for _, n := range notes {
		data = append(data, Data{n, idx})
		idx--
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func showAddFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "form", nil)
}

func closeAddFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "share_button", nil)
}

func addNoteHandler(w http.ResponseWriter, r *http.Request) {
	author := r.FormValue("author")
	message := r.FormValue("message")

	if author == "" {
		author = "stranger"
	}

	resp, err := http.Get("http://localhost:3000/notes/count")
	if err != nil {
		log.Fatalf("Failed to fetch notes count: %s", err)
	}
	countStr, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	count, _ := strconv.Atoi(string(countStr))

	type Data struct {
		Note
		Idx int
	}
	note := Note{}

	resp, err = http.Post("http://localhost:3000/notes", "application/json", strings.NewReader(fmt.Sprintf(`{"author": %q, "message": %q}`, author, message)))
	if err != nil {
		log.Fatalf("Failed to post a note: %s", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&note)
	if err != nil {
		log.Fatalf("Failed to decode response: %s", err)
	}

	data := Data{note, count + 1}

	tmpl.ExecuteTemplate(w, "share_button", nil)
	tmpl.ExecuteTemplate(w, "note", map[string]any{"Note": data, "NoteCreated": true})
}
