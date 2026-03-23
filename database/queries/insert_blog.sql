-- name: InsertBlog :one
insert into blogs(content, user_id) values(?, ?)

returning *;
