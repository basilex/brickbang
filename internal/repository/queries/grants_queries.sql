-- name: CreateGrant :one
INSERT INTO grants (code, description) VALUES (@code, @description)
RETURNING *;

-- name: CountGrants :one
SELECT COUNT(*) FROM grants;

-- name: GetGrantByID :one
SELECT * FROM grants WHERE pid = @pid;

-- name: GetGrantByCode :one
SELECT * FROM grants WHERE code = @code;

-- name: ListGrants :many
SELECT *
  FROM grants
 ORDER BY
    CASE WHEN @sql_order = 'asc' THEN code END ASC,
    CASE WHEN @sql_order = 'desc' THEN code END DESC,
      code ASC
 LIMIT @sql_limit OFFSET @sql_offset;

-- name: UpdateGrantByID :one
UPDATE grants
   SET code = @code, description = @description
 WHERE pid = @pid RETURNING *;

-- name: DeleteGrantByID :one
DELETE FROM grants WHERE pid = @pid RETURNING *;
