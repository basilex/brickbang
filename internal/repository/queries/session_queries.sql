-- name: CreateSession :one
INSERT INTO session (
  user_id, access_token, refresh_token, access_exp, refresh_exp, access_jti, refresh_jti, ip_address, user_agent
) VALUES (@user_id,@access_token,@refresh_token,@access_exp,@refresh_exp,@access_jti,@refresh_jti,@ip_address,@user_agent)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM session WHERE id = @id;

-- name: ListSessionsByUserID :many
SELECT * FROM session WHERE user_id = @user_id ORDER BY created_at DESC;

-- name: UpdateAccessTokenByID :one
UPDATE session
   SET access_token = @access_token, access_jti = @access_jti, access_exp = @access_exp
 WHERE id = @id RETURNING *;

-- name: RotateTokensByID :one
UPDATE session
SET access_token = @access_token, access_jti = @access_jti, access_exp = @access_exp,
    refresh_token = @refresh_token, refresh_jti = @refresh_jti, refresh_exp = @refresh_exp,
    updated_at = timezone('utc', now())
WHERE id = @id
RETURNING *;

-- name: ExpireAccessTokenByJTI :one
UPDATE session
   SET access_status = 'expired', updated_at = timezone('utc', now())
WHERE access_jti = @access_jti
RETURNING *;

-- name: RevokeAccessTokenByJTI :one
UPDATE session
   SET access_status = 'revoked', updated_at = timezone('utc', now())
 WHERE access_jti = @access_jti RETURNING *;

-- name: RevokeSessionByID :one
UPDATE session
   SET access_status = 'revoked', refresh_status = 'revoked',
    updated_at = timezone('utc', now())
WHERE id = @id
RETURNING *;
