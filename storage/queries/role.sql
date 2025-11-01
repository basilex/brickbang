-- name: RoleNew :one
insert into role(name) values(@name) returning *;

-- name: RoleCount :one
select count(*) from role;

-- name: RoleSelect :many
select *
  from role r
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: RoleSelectByID :one
select * from role r where r.id = @id;

-- name: RoleSelectByName :one
select * from role r where r.name = @name;

-- name: RoleUpdateByID :one
update role
   set name = @name where id = @id returning *;

-- name: RoleDeleteByID :one
delete from role where id = @id returning id;