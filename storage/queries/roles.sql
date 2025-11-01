-- name: RolesNew :one
insert into roles(name) values(@name) returning *;

-- name: RolesCount :one
select count(*) from roles;

-- name: RolesSelect :many
select *
  from roles r
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: RolesSelectByID :one
select * from roles r where r.id = @id;

-- name: RolesSelectByName :one
select * from roles r where r.name = @name;

-- name: RolesUpdateByID :one
update roles
   set name = @name where id = @id returning *;

-- name: RolesDeleteByID :one
delete from roles where id = @id returning id;
