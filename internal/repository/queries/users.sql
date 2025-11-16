-- Users

-- Create
-- name: CreateUser :one
INSERT INTO users (username, password, is_checked)
VALUES ($1, $2, $3)
RETURNING *;

-- Get by ID
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- Get by username
-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- List all users
-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- Update user
-- name: UpdateUserByID :one
UPDATE users
SET username = $2, password = $3, is_checked = $4, is_blocked = $5
WHERE id = $1
RETURNING *;

-- Update block status
-- name: UpdateUserIsBlockedByID :one
UPDATE users
SET is_blocked = $2,
    blocked_at = CASE WHEN $2 THEN NOW() ELSE '1000-01-01'::timestamp END,
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, is_blocked, is_checked;

-- Delete user
-- name: DeleteUserByID :one
DELETE FROM users WHERE id = $1
RETURNING *;
