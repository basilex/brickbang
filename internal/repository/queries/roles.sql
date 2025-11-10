-- name: RoleCreate :one
insert into roles(name) values(@name) returning *;

-- name: RolesCount :one
select count(*) from roles;

-- name: RolesList :many
select *
from roles r
order by @sql_order::text
limit @sql_limit offset @sql_offset;

-- name: RoleGetByID :one
select * from roles r where r.id = @id;

-- name: RoleGetByName :one
select * from roles r where r.name = @name;

-- name: RoleUpdateByID :one
update roles
set name = @name
where id = @id
returning *;

-- name: RoleDeleteByID :one
delete from roles where id = @id returning id;
