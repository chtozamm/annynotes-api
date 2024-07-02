package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	// Check if Content-Type header is present
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0])) // normalize header by stripping whitespace and any additional parameters
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	// Restrict request body size
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup decoder and disallow unknown fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		// Catch any syntax errors in the JSON
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// In some circumstances Decode() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch any type errors
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// An io.EOF error is returned by Decode() if the request body is
		// empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// Catch the error caused by the request body being too large
		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("Request body must not be larger than 1MB (got %d bytes)", maxBytesError.Limit)
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		// Otherwise default to logging the error and sending a 500
		// Internal Server Error response.
		default:
			return err
		}
	}

	// Call decode again, using a pointer to an empty anonymous struct as
	// the destination. If the request body only contained a single JSON
	// object this will return an io.EOF error. So if we get anything else,
	// we know that there is additional data in the request body.
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func setupDB(conn *sql.DB) error {
	_, err := conn.Exec(`
CREATE TABLE IF NOT EXISTS notes (
  id TEXT NOT NULL PRIMARY KEY,
  author TEXT NOT NULL,
  message TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  user_id TEXT NOT NULL DEFAULT "mmq9tll8as33g64",
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_note_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
  UPDATE notes
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;

CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%fZ', 'now')),
  verified INTEGER NOT NULL DEFAULT 0
);

CREATE TRIGGER IF NOT EXISTS update_user_timestamp
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
  UPDATE users
  SET updated_at = strftime('%Y-%m-%d %H:%M:%fZ', 'now')
  WHERE id = NEW.id;
END;
`)
	if err != nil {
		return err
	}
	return nil
}
