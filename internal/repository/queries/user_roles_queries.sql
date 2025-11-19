-- name: CreateUserRole :one
INSERT INTO user_roles (user_id, role_id) VALUES (@user_id, @role_id) RETURNING *;

-- name: GetUserRoleByID :one
SELECT * FROM user_roles WHERE id = @id;
-- name: ListUserRolesByUserID :many
SELECT * FROM user_roles WHERE user_id = @user_id;

-- name: ListUserRolesByRoleID :many
SELECT * FROM user_roles WHERE role_id = @role_id;

-- name: DeleteUserRoleByID :one
DELETE FROM user_roles WHERE id = @id RETURNING *;
