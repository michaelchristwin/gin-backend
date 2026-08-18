-- name: CreateUser :one
INSERT INTO users (name, email) VALUES (?, ?) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id LIMIT ? OFFSET ?;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(sqlc.narg(name), name),
    email = COALESCE(sqlc.narg(email), email)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = ?;

