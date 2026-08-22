-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash
) VALUES (?, ?)
RETURNING *;

-- name: GetUser :one
SELECT
    id,
    email,
    created_at
FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;

-- name: ListUsers :many
SELECT
    id,
    email,
    created_at
FROM users
ORDER BY id
LIMIT ? OFFSET ?;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(sqlc.narg(email), email),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetTotalUsers :one
SELECT COUNT(*)
FROM users;

-- name: CreateSession :one
INSERT INTO sessions (
    id,
    user_id,
    expires_at
) VALUES (?, ?, ?)
RETURNING *;

-- name: GetSession :one
SELECT *
FROM sessions
WHERE id = ?
LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= CURRENT_TIMESTAMP;

-- name: GetSessionWithUser :one
SELECT
    s.id AS session_id,
    s.user_id,
    s.expires_at,
    u.email,
    u.password_hash,
    u.created_at
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.id = ?
LIMIT 1;