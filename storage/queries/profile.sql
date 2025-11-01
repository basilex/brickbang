-- name: ProfileNew :one
insert into profile(
    user_id, firstname, lastname
) values(
    @user_id, @firstname, @lastname
) returning *;

-- name: ProfileCount :one
select count(*) from profile;

-- name: ProfileSelect :many
select *
  from profile p
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: ProfileSelectByID :one
select * from profile p where p.id = @id;

-- name: ProfileSelectByUserID :one
select * from profile p where p.user_id = @id;

-- name: ProfileUpdateCommonById :one
update profile
   set firstname = @firstname, lastname = @lastname, gender = @gender, birthday = @birthday
 where id = @id returning *;

-- name: ProfileUpdate2FAById :one
update profile
   set enable_2fa = @enable_2fa, secret_2fa = @secret_2fa
 where id = @id returning *;

-- name: ProfileUpdateAvatarUrlById :one
update profile
   set avatar_url = @avatar_url
 where id = @id returning *;

-- name: ProfileIsExists :one
select case when exists (
    select * from profile p where p.user_id = @user_id
) then cast(1 as bit) else cast(0 as bit) end;
