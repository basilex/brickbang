-- name: ContactNew :one
insert into Contact(
    user_id, class, content
) values(
    @user_id, @class, @content
) returning *;

-- name: ContactCount :one
select count(*) from contact;

-- name: ContactCountByUserID :one
select count(*) from contact where user_id = @user_id;

-- name: ContactSelect :many
select *
  from contact
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: ContactSelectByID :one
select * from contact c where c.id = @id;

-- name: ContactSelectByUserID :many
select * from contact c where c.user_id = @user_id order by c.class;

-- name: ContactSelectByUserIDClass :many
select * from contact c
 where c.user_id = @user_id
   and c.class = @class
   and c.content = @content;

-- name: ContactSelectUserByClass :one
select u.id, c.content as email, u.username, u.is_blocked, u.blocked_at, u.is_checked, u.checked_at, u.visited_at, u.created_at, u.updated_at
  from users u
  join contact c on u.id = c.user_id
where c.class = @class
  and c.content = @content;

-- name: ContactUpdateByID :one
update contact
   set class = @class, content = @content
 where id = @id returning *;

-- name: ContactDeleteByID :one
delete from contact c where c.id = @id returning id;