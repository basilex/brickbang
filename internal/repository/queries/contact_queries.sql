-- name: CreateContact :one
INSERT INTO contact (user_id, class, content)
VALUES (@user_id,@class,@content) RETURNING *;

-- name: GetContactByID :one
SELECT * FROM contact WHERE id = @id;

-- name: ListContactsByUserID :many
SELECT * FROM contact WHERE user_id = @user_id ORDER BY created_at DESC;

-- Update
-- name: UpdateContactByID :one
UPDATE contact
   SET class = @class, content = @content
 WHERE id = @id RETURNING *;

-- name: DeleteContactByID :one
DELETE FROM contact WHERE id = @id RETURNING *;
