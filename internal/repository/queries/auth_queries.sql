-- name: AuthCreateUser :one
INSERT INTO users (username, password, is_checked)
VALUES (@username, @password, @is_checked)
RETURNING id, username, password, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: AuthSelectUserByID :one
SELECT id, username, password, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at
  FROM users
 WHERE id = @id;

-- name: AuthSelectUserCredentials :one
SELECT id, username, password, is_blocked, is_checked, blocked_at, checked_at
  FROM users
 WHERE username = @username;

-- name: AuthUpdateVisitedAt :one
UPDATE users
   SET visited_at = timezone('utc', now())
 WHERE id = @id RETURNING id, username, visited_at, created_at, updated_at;

-- name: AuthBlockUserByID :one
UPDATE users
   SET is_blocked = true, blocked_at = timezone('utc', now())
 WHERE id = @id RETURNING id, username, is_blocked, blocked_at;

-- name: AuthGetUserGrants :many
SELECT g.code
  FROM user_roles ur
  JOIN role_grants rg ON rg.role_id = ur.role_id
  JOIN grants g       ON g.id = rg.grant_id
 WHERE ur.user_id = @user_id;
