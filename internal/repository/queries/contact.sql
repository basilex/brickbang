-- name: ListContactsByUserID :many
select *
from contact
where user_id = @user_id
order by class, content;

-- name: GetContactByID :one
select *
from contact
where id = @id;

-- name: CreateContact :one
insert into contact(user_id, class, content)
values(@user_id, @class, @content)
returning *;

-- name: UpdateContact :one
update contact
set class = @class,
    content = @content
where id = @id
returning *;

-- name: DeleteContactByID :one
delete from contact where id = @id returning id;
