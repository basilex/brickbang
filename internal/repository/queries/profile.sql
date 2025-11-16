-- Profile

-- Create
-- name: CreateProfile :one
INSERT INTO profile (user_id, firstname, lastname, gender, birthday, avatar_url, enable_2fa, secret_2fa)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- Get by ID
-- name: GetProfileByID :one
SELECT * FROM profile WHERE id = $1;

-- Get by UserID
-- name: GetProfileByUserID :one
SELECT * FROM profile WHERE user_id = $1;

-- List all
-- name: ListProfiles :many
SELECT * FROM profile ORDER BY created_at DESC;

-- Update
-- name: UpdateProfileByID :one
UPDATE profile
SET firstname=$2, lastname=$3, gender=$4, birthday=$5, avatar_url=$6, enable_2fa=$7, secret_2fa=$8
WHERE id=$1
RETURNING *;

-- Delete
-- name: DeleteProfileByID :one
DELETE FROM profile WHERE id=$1
RETURNING *;
