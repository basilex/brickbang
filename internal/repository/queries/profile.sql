-- name: GetProfileByUserID :one
select *
from profile
where user_id = @user_id;

-- name: UpdateProfileByUserID :one
update profile
set firstname = @firstname,
    lastname = @lastname,
    gender = @gender,
    birthday = @birthday,
    avatar_url = @avatar_url,
    enable_2fa = @enable_2fa,
    secret_2fa = @secret_2fa
where user_id = @user_id
returning *;
