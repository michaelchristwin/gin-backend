-- name: CreateUser :one
INSERT INTO users (name, email) VALUES (?, ?) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id LIMIT ? OFFSET ?;

-- name: UpdateUser :one
UPDATE users 
SET name = ?, email = ? 
WHERE id = ? 
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users 
WHERE id = ? 
RETURNING *;