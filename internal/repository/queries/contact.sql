-- Contact

-- Create
-- name: CreateContact :one
INSERT INTO contact (user_id, class, content)
VALUES ($1,$2,$3)
RETURNING *;

-- Get by ID
-- name: GetContactByID :one
SELECT * FROM contact WHERE id = $1;

-- List by UserID
-- name: ListContactsByUserID :many
SELECT * FROM contact WHERE user_id = $1 ORDER BY created_at DESC;

-- Update
-- name: UpdateContactByID :one
UPDATE contact
SET class=$2, content=$3
WHERE id=$1
RETURNING *;

-- Delete
-- name: DeleteContactByID :one
DELETE FROM contact WHERE id=$1
RETURNING *;
