-- name: CreateUser :one
INSERT INTO users (username, password, is_checked)
VALUES (@username, @password, @is_checked) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = @id;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = @username;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: UpdateUserByID :one
UPDATE users
   SET username = @username, password = @password, is_checked = @is_checked, is_blocked = @is_blocked
 WHERE id = @id RETURNING *;

-- name: UpdateVisitedAtByUserID :one
UPDATE users SET visited_at = @visited_at WHERE id = @id RETURNING *;

-- name: UpdateUserIsBlockedByID :one
UPDATE users
   SET is_blocked = @is_blocked,
       blocked_at = CASE WHEN @is_blocked THEN NOW() ELSE '1000-01-01'::timestamp END,
       updated_at = NOW()
 WHERE id = @id RETURNING id, username, is_blocked, is_checked;

-- name: DeleteUserByID :one
DELETE FROM users WHERE id = @id RETURNING *;
