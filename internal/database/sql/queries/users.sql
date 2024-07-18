-- name: CreateUser :one
INSERT INTO users (id, email, name, username, password)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;