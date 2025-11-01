-- name: UserNew :one
insert into users(
    username, password
) values(
    @username, @password
) returning id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: UserCount :one
select count(*) from users;

-- name: UserSelect :many
select id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at
  from users u
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: UserSelectByID :one
select id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at
  from users u
 where u.id = @id;

-- name: UserSelectByUsername :one
select id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at
  from users u
 where u.username = @username;

-- name: UserUpdateCredentialsByID :one
update users
   set username = @username, password = @password
 where id = @id
       returning id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: UserUpdateIsBlockedByID :one
update users
   set is_blocked = @is_blocked,
       blocked_at = case
                        when @is_blocked then timezone('utc', now())
                        else blocket_at
                    end
 where id = @id
       returning id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: UserUpdateIsCheckedByID :one
update users
   set is_checked = @is_checked,
       checked_at = case
                        when @is_checked then timezone('utc', now())
                        else checked_at
                    end
 where id = @id
       returning id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: UserUpdateVisitedAtByID :one
update users
   set visited_at = timezone('utc', now())
 where id = @id
       returning id, username, is_blocked, is_checked, blocked_at, checked_at, visited_at, created_at, updated_at;

-- name: UserDeleteByID :one
delete from users where id = @id returning id;