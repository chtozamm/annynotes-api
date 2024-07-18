package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/chtozamm/annynotes-go/internal/auth"
	"github.com/chtozamm/annynotes-go/internal/database"
	"github.com/chtozamm/annynotes-go/internal/utils"
	"github.com/mattn/go-sqlite3"
)

func (app *application) getNotesHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(notes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	payload, err := json.Marshal(&struct {
		Total int             `json:"total"`
		Notes []database.Note `json:"notes"`
	}{
		Total: len(notes),
		Notes: notes,
	})
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
		return
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

	payload, err := json.Marshal(&struct {
		Total int             `json:"total"`
		Notes []database.Note `json:"notes"`
	}{
		Total: len(notes),
		Notes: notes,
	})
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
		return
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

	if note.Author == "" {
		note.Author = "stranger"
	}

	note.ID = utils.GenerateUniqueId()

	if !utils.ValidateId(note.ID) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
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
		return
	}

	if !utils.ValidateId(id) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	note, err := app.DB.FetchNoteByID(r.Context(), id)
	if err != nil {
		log.Printf("Attempt to delete a non-existing note with the ID: %q", id)
		http.Error(w, "Note does not exist", http.StatusNotFound)
		return
	}
	if note.UserID != user.ID {
		log.Printf("Unauthorized attempt to delete a note %q by user %q", note.ID, user.ID)
		http.Error(w, "Note belongs to another user", http.StatusUnauthorized)
		return
	}

	err = app.DB.DeleteNote(r.Context(), id)
	if err != nil {
		log.Printf("Failed to delete a note with the ID %q: %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
		return
	}

	if !utils.ValidateId(id) {
		msg := "Invalid note ID"
		log.Printf(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	note, err := app.DB.FetchNoteByID(r.Context(), id)
	if err != nil {
		log.Printf("Attempt to update a non-existing note with the ID: %q", id)
		http.Error(w, "Note does not exist", http.StatusNotFound)
		return
	}
	if note.UserID != user.ID {
		log.Printf("Unauthorized attempt to update a note %q by user %q", note.ID, user.ID)
		http.Error(w, "Note belongs to another user", http.StatusUnauthorized)
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

	if newNote.Message == "" {
		msg := "Message field cannot be empty"
		log.Print("Tried to update a note without a message")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if newNote.Author == "" {
		newNote.Author = note.Author
	}

	updatedNote, err := app.DB.UpdateNote(r.Context(), database.UpdateNoteParams{
		ID:      id,
		Author:  newNote.Author,
		Message: newNote.Message,
	})
	if err != nil {
		log.Printf("Failed to update a note with the ID %q: %s", id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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

	user.ID = utils.GenerateUniqueId()

	switch {
	case user.Email == "":
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	case user.Password == "":
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	case user.Username == "":
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	case user.Name == "":
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if len(user.Password) < 4 {
		log.Print("Provided password is too short")
		http.Error(w, "Password must contain at least 4 characters", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		log.Printf("Failed to hash password: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Validate email
	email, err := mail.ParseAddress(user.Email)
	if err != nil {
		log.Printf("Failed to parse an email: %s", err)
		http.Error(w, "Email is not valid", http.StatusBadRequest)
		return
	}

	newUser, err := app.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       user.ID,
		Email:    email.Address,
		Name:     user.Name,
		Username: user.Username,
		Password: hashedPassword,
	})
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			log.Print(sqliteErr.Code)
			log.Print(sqliteErr.Error())
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				log.Printf("Attempt to create a new user with email that already exists: %s", err)
				http.Error(w, "Email is already in use", http.StatusConflict)
				return
			case sqlite3.ErrConstraintCheck:
				if strings.Contains(sqliteErr.Error(), "length(username)") {
					http.Error(w, "Username must be between 2 and 20 characters long", http.StatusBadRequest)
					return
				}
				if strings.Contains(sqliteErr.Error(), "length(name)") {
					http.Error(w, "Name must be between 2 and 20 characters long", http.StatusBadRequest)
					return
				}
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			default:
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
		log.Printf("Failed to create new user: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Printf("New user created with the ID %q", user.ID)
	auth.RespondWithJWT(w, newUser.ID, newUser.Email)
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

	storedUser, err := app.DB.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		// Send 404 if user doesn't exist
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User does not exist", http.StatusNotFound)
			return
		}
		// Send 500 for any other errors
		log.Printf("User authentication fail: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check password against the hash
	if !auth.CheckPassword(storedUser.Password, user.Password) {
		log.Print("Attempt to login with incorrect password")
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	auth.RespondWithJWT(w, storedUser.ID, storedUser.Email)
}
