package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/chtozamm/annynotes-go/internal/auth"
	"github.com/chtozamm/annynotes-go/internal/database"
	"github.com/chtozamm/annynotes-go/internal/utils"
	"github.com/mattn/go-sqlite3"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	var notes []database.Note
	var err error

	// Fetch notes ordered according to the URL query
	sortQuery := r.URL.Query().Get("sort")
	switch strings.ToLower(sortQuery) {
	case "desc":
		notes, err = app.DB.FetchNotesDESC(r.Context())
	default:
		notes, err = app.DB.FetchNotes(r.Context())
	}

	if err != nil {
		log.Printf("Failed to fetch notes: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Handle no content response
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"notes": []}`))
		return
	}

	payload, err := json.Marshal(&struct {
		Notes []database.Note `json:"notes"`
	}{Notes: notes})
	if err != nil {
		log.Printf("Failed to marshal notes: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) getNotesFromAuthorHandler(w http.ResponseWriter, r *http.Request) {
	author := r.PathValue("author")
	if author == "" {
		msg := "Author was not provided"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	// Normalize author name
	author = utils.NormalizeName(author)

	var notes []database.Note
	var err error

	// Fetch notes ordered according to the URL query
	sortQuery := r.URL.Query().Get("sort")
	switch strings.ToLower(sortQuery) {
	case "desc":
		notes, err = app.DB.FetchNotesFromAuthorDESC(r.Context(), author)
	default:
		notes, err = app.DB.FetchNotesFromAuthor(r.Context(), author)
	}

	// Handle no content response
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"notes": []}`))
		return
	}

	if err != nil {
		log.Printf("Failed to fetch notes from author: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Handle no content response
	if len(notes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write([]byte(`{"notes": []}`))
		return
	}

	payload, err := json.Marshal(&notes)
	if err != nil {
		log.Printf("Failed to marshal notes: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) getNoteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		msg := "Note ID was not provided"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	if !utils.ValidateId(id) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	note, err := app.DB.FetchNoteByID(r.Context(), id)
	if err != nil {
		log.Printf("Failed to fetch note with the ID %q: %s", id, err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	payload, err := json.Marshal(&note)
	if err != nil {
		log.Printf("Failed to marshal note: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) createNoteHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	var note database.Note

	err := decodeJSONBody(w, r, &note)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if note.Message == "" {
		msg := "Message field cannot be empty"
		log.Print("Tried to create a new note without a message")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if note.ID == "" {
		note.ID = utils.GenerateUniqueId()
	}

	if !utils.ValidateId(note.ID) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	newNote, err := app.DB.CreateNote(r.Context(), database.CreateNoteParams{
		ID:       note.ID,
		Author:   note.Author,
		Message:  note.Message,
		UserID:   user.ID,
		Verified: user.Verified,
	})
	if err != nil {
		log.Printf("Failed to create a new note: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(&newNote)
	if err != nil {
		log.Printf("Failed to marshal note: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Printf("New note created with the ID %q", note.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) deleteNoteHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	id := r.PathValue("id")

	if id == "" {
		msg := "Note ID was not provided"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	if !utils.ValidateId(id) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	note, err := app.DB.FetchNoteByID(r.Context(), id)
	if err != nil {
		log.Printf("Attempt to delete a non-existing note with the ID: %q", id)
		http.Error(w, "Note does not exist", http.StatusNotFound)
		return
	}
	if note.UserID != user.ID {
		log.Printf("Unauthorized attempt to delete a note %q by user %q", note.ID, user.ID)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = app.DB.DeleteNote(r.Context(), id)
	if err != nil {
		log.Printf("Failed to delete a note with the ID %q: %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	log.Printf("Delete a note with the ID %q", id)
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateNoteHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	id := r.PathValue("id")

	if id == "" {
		msg := "Note ID was not provided"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	if !utils.ValidateId(id) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	note, err := app.DB.FetchNoteByID(r.Context(), id)
	if err != nil {
		log.Printf("Attempt to update a non-existing note with the ID: %q", id)
		http.Error(w, "Note does not exist", http.StatusNotFound)
		return
	}
	if note.UserID != user.ID {
		log.Printf("Unauthorized attempt to update a note %q by user %q", note.ID, user.ID)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var newNote database.Note

	err = decodeJSONBody(w, r, &newNote)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	updatedNote, err := app.DB.UpdateNote(r.Context(), database.UpdateNoteParams{
		ID:      id,
		Author:  newNote.Author,
		Message: newNote.Message,
	})
	if err != nil {
		log.Printf("Failed to update a note with the ID %q: %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	log.Printf("Update a note with the ID %q", id)

	payload, err := json.Marshal(&updatedNote)
	if err != nil {
		log.Printf("Failed to marshal note: %s", err)
		http.Error(w, "Note was updated successfully, but the server couldn't respond with the updated note.", http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {

	var user database.User

	err := decodeJSONBody(w, r, &user)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if user.ID == "" {
		user.ID = utils.GenerateUniqueId()
	}

	if !utils.ValidateId(user.ID) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}

	newUser, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint {
				log.Printf("Attempt to create a new user with email that already exists: %s", err)
				http.Error(w, "Email is already in use", http.StatusConflict)
				return
			}
		}
		log.Printf("Failed to create new user: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateJWT(newUser.ID, newUser.Email, newUser.Password)
	if err != nil {
		log.Printf("Failed to generate JWT: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(&struct {
		Token string `json:"token"`
	}{Token: token})
	if err != nil {
		log.Printf("Failed to marshal token: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Printf("New user created with the ID %q", user.ID)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {

	var user database.User

	err := decodeJSONBody(w, r, &user)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if user.Email == "" || user.Password == "" {
		msg := "Malformed request: expected payload to have email and password fields"
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, err := app.DB.GetUserIDByCredentials(r.Context(), database.GetUserIDByCredentialsParams{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		log.Printf("User authentication fail: %s", err)
		http.Error(w, "Wrong email or password", http.StatusNotFound)
		return
	}
	user.ID = id

	token, err := auth.GenerateJWT(user.ID, user.Email, user.Password)
	if err != nil {
		log.Printf("Failed to generate JWT: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(&struct {
		Token string `json:"token"`
	}{Token: token})
	if err != nil {
		log.Printf("Failed to marshal token: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}
