-- name: InsertBlog :one
insert into blogs(content, title, category, tags, createdAt, updatedAt) values(?, ?, ?, ?, ?, ?)

returning *;
