-- UserRoles

-- Create
-- name: CreateUserRole :one
INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)
RETURNING *;

-- Get by ID
-- name: GetUserRoleByID :one
SELECT * FROM user_roles WHERE id = $1;

-- List by UserID
-- name: ListUserRolesByUserID :many
SELECT * FROM user_roles WHERE user_id = $1;

-- List by RoleID
-- name: ListUserRolesByRoleID :many
SELECT * FROM user_roles WHERE role_id = $1;

-- Delete
-- name: DeleteUserRoleByID :one
DELETE FROM user_roles WHERE id = $1
RETURNING *;
