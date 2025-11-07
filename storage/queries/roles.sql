-- name: CreateRole :one
insert into roles(name) values(@name) returning *;

-- name: CountRoles :one
select count(*) from roles;

-- name: ListRoles :many
select *
  from roles r
 order by @sql_order::text
 limit @sql_limit offset @sql_offset;

-- name: GetRoleByID :one
select * from roles r where r.id = @id;

-- name: GetRoleByName :one
select * from roles r where r.name = @name;

-- name: UpdateRoleByID :one
update roles
   set name = @name where id = @id returning *;

-- name: DeleteRoleByID :one
delete from roles where id = @id returning id;
