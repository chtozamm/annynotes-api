-- name: CreateNote :one
INSERT INTO notes (id, author, message, user_id, verified) 
VALUES (?, ?, ?, ?, ?) 
RETURNING *;

-- name: UpdateNote :one
UPDATE notes SET author = ?, message = ? 
WHERE id = ?
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = ?;

-- name: FetchNotes :many
SELECT * FROM notes
ORDER BY created_at ASC;

-- name: FetchNotesDESC :many
SELECT * FROM notes
ORDER BY created_at DESC;

-- name: FetchNotesFromAuthor :many
SELECT * FROM notes
WHERE author = ?
ORDER BY created_at ASC;

-- name: FetchNotesFromAuthorDESC :many
SELECT * FROM notes
WHERE author = ?
ORDER BY created_at DESC;

-- name: FetchNoteByID :one
SELECT * FROM notes WHERE id = ?;
