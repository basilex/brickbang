-- name: CreateRole :one
INSERT INTO roles (name, description) VALUES (@name, @description)
RETURNING *;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = @id;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = @name;

-- name: ListRoles :many
SELECT *
  FROM roles
 ORDER BY
    CASE WHEN @sql_order = 'asc' THEN name END ASC,
    CASE WHEN @sql_order = 'desc' THEN name END DESC,
      name ASC
 LIMIT @sql_limit OFFSET @sql_offset;

-- name: UpdateRoleByID :one
UPDATE roles
   SET name = @name, description = @description
 WHERE id = @id RETURNING *;

-- name: DeleteRoleByID :one
DELETE FROM roles WHERE id = @id RETURNING *;
