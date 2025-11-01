-- name: AuthSelectUserCredentials :one
select id, username, password, is_blocked, is_checked, blocked_at, checked_at
  from users u
 where u.username = @username;

-- name: AuthUpdateVisitedAt :one
update users set visited_at = now() where id = @id
       returning id, username, checked_at, visited_at, created_at, updated_at;
