-- Session

-- name: CreateSession :one
INSERT INTO session (
  user_id, access_token, refresh_token, access_exp, refresh_exp, access_jti, refresh_jti, ip_address, user_agent
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM session WHERE id = $1;

-- name: ListSessionsByUserID :many
SELECT * FROM session WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateAccessTokenByID :one
UPDATE session
   SET access_token = $2, access_jti = $3, access_exp = $4
 WHERE id = $1 RETURNING *;

-- name: RotateTokensByID :one
UPDATE session
SET access_token = $2, access_jti = $3, access_exp = $4,
    refresh_token = $5, refresh_jti = $6, refresh_exp = $7,
    updated_at = timezone('utc', now())
WHERE id = $1
RETURNING *;

-- name: ExpireAccessTokenByJTI :one
UPDATE session
   SET access_status = 'expired', updated_at = timezone('utc', now())
WHERE access_jti = $1
RETURNING *;

-- name: RevokeAccessTokenByJTI :one
UPDATE session
   SET access_status = 'revoked', updated_at = timezone('utc', now())
 WHERE access_jti = $1 RETURNING *;
