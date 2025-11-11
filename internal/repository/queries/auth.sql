-- USERS AUTH MANAGEMENT

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
set visited_at = timezone('utc', now())
where id = @id
returning id, username, visited_at, created_at, updated_at;

-- name: AuthBlockUserByID :one
update users
set is_blocked = true,
    blocked_at = timezone('utc', now())
where id = @id
returning id, username, is_blocked, blocked_at;

-- SESSION MANAGEMENT

-- name: AuthCreateSession :one
insert into session (
    user_id,
    access_token,
    refresh_token,
    access_exp,
    refresh_exp,
    ip_address,
    user_agent
) values (
    @user_id, @access_token, @refresh_token, @access_exp, @refresh_exp, @ip_address, @user_agent
)
returning
    id, user_id, access_token, refresh_token, access_exp, refresh_exp,
    access_status, refresh_status, ip_address, user_agent, created_at, updated_at;

-- name: AuthSelectSessionByAccessToken :one
select *
from session
where access_token = @access_token
  and access_status = 'valid';

-- name: AuthSelectSessionByRefreshToken :one
select *
from session
where refresh_token = @refresh_token
  and refresh_status = 'valid';

-- name: AuthRevokeAccessSessionByID :one
update session
set access_status = 'revoked'
where id = @id
returning id, user_id, access_status;

-- name: AuthRevokeRefreshSessionByID :one
update session
set refresh_status = 'revoked'
where id = @id
returning id, user_id, refresh_status;

-- name: AuthExpireAccessTokenByID :one
update session
set access_status = 'expired'
where id = @id
returning id, access_status;

-- name: AuthExpireRefreshTokenByID :one
update session
set refresh_status = 'expired'
where id = @id
returning id, refresh_status;

-- name: AuthUpdateAccessTokenByID :one
update session
set access_token = @access_token,
    access_exp = @access_exp,
    access_status = 'valid'
where id = @id
returning id, access_token, access_exp, access_status;

-- name: AuthUpdateRefreshTokenByID :one
update session
set refresh_token = @refresh_token,
    refresh_exp = @refresh_exp,
    refresh_status = 'valid'
where id = @id
returning id, refresh_token, refresh_exp, refresh_status;

-- name: AuthListSessionsByUserID :many
select *
from session
where user_id = @user_id
order by created_at desc;

-- CLEANUP / ADMIN

-- name: AuthDeleteSessionsByUserID :exec
delete from session
where user_id = @user_id;

-- name: AuthDeleteExpiredSessions :exec
delete from session
where (access_exp < timezone('utc', now()) and access_status != 'revoked')
   or (refresh_exp < timezone('utc', now()) and refresh_status != 'revoked');
