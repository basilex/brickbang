-- name: CreateProfile :one
INSERT INTO profile (user_id, firstname, lastname, gender, birthday, avatar_url, enable_2fa, secret_2fa)
VALUES (@user_id,@firstname,@lastname,@gender,@birthday,@avatar_url,@enable_2fa,@secret_2fa)
RETURNING *;

-- name: GetProfileByID :one
SELECT * FROM profile WHERE id = @id;

-- name: GetProfileByUserID :one
SELECT * FROM profile WHERE user_id = @user_id;

-- name: ListProfiles :many
SELECT * FROM profile ORDER BY created_at DESC;

-- name: UpdateProfileByID :one
UPDATE profile
   SET firstname = @firstname, lastname = @lastname, gender = @gender, birthday = @birthday,
       avatar_url = @avatar_url, enable_2fa = @enable_2fa, secret_2fa = @secret_2fa
 WHERE id = @id RETURNING *;

-- name: DeleteProfileByID :one
DELETE FROM profile WHERE id = @id RETURNING *;
