-- Users

-- name: AuthSelectUserByID :one
select id, username, password, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at
from users
where id = @id;

-- name: AuthSelectUserCredentials :one
select id, username, password, is_blocked, is_checked, blocked_at, checked_at
from users
where username = @username;

-- name: AuthCreateUser :one
insert into users (username, password, is_checked)
values (@username, @password, @is_checked)
returning id, username, password, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: AuthUpdateVisitedAt :one
update users
set visited_at = now()
where id = @id
returning id, username, checked_at, visited_at, created_at, updated_at;

