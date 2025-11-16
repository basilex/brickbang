-- Roles

-- Create
-- name: CreateRole :one
INSERT INTO roles (name) VALUES ($1)
RETURNING *;

-- Count
-- name: CountRoles :one
SELECT COUNT(*) FROM roles;

-- Get by ID
-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;

-- Get by name
-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;

-- List with pagination
-- name: ListRoles :many
SELECT *
FROM roles
ORDER BY
    CASE WHEN @sqlorder = 'asc' THEN name END ASC,
    CASE WHEN @sqlorder = 'desc' THEN name END DESC
LIMIT @sqllimit OFFSET @sqloffset;

-- Update
-- name: UpdateRoleByID :one
UPDATE roles
SET name = $2
WHERE id = $1
RETURNING *;

-- Delete
-- name: DeleteRoleByID :one
DELETE FROM roles WHERE id = $1
RETURNING *;
