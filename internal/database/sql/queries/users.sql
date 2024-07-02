-- name: CreateUser :one
INSERT INTO users (id, email, password)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserIDByCredentials :one
SELECT id FROM users WHERE email = ? AND password = ?;