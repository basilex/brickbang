-- name: CreateSession :one
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

-- name: CountSessions :one
select count(*) from session;

-- name: ListSessionsByUserID :many
select *
from session
where user_id = @user_id
order by created_at desc;

-- name: GetSessionByID :one
select *
from session
where id = @id;

-- name: GetSessionByAccessToken :one
select *
from session
where access_token = @access_token
  and access_status = 'valid';

-- name: GetSessionByRefreshToken :one
select *
from session
where refresh_token = @refresh_token
  and refresh_status = 'valid';

-- name: UpdateAccessTokenByID :one
update session
set access_token = @access_token,
    access_exp = @access_exp,
    access_status = 'valid'
where id = @id
returning id, access_token, access_exp, access_status;

-- name: UpdateRefreshTokenByID :one
update session
set refresh_token = @refresh_token,
    refresh_exp = @refresh_exp,
    refresh_status = 'valid'
where id = @id
returning id, refresh_token, refresh_exp, refresh_status;

-- name: RevokeAccessTokenByID :one
update session
set access_status = 'revoked'
where id = @id
returning id, user_id, access_status;

-- name: RevokeRefreshTokenByID :one
update session
set refresh_status = 'revoked'
where id = @id
returning id, user_id, refresh_status;

-- name: ExpireAccessTokenByID :one
update session
set access_status = 'expired'
where id = @id
returning id, access_status;

-- name: ExpireRefreshTokenByID :one
update session
set refresh_status = 'expired'
where id = @id
returning id, refresh_status;
