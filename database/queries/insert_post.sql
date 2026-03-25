-- name: InsertPost :one
insert into posts(content, title, category, tags, createdAt, updatedAt) values(?, ?, ?, ?, ?, ?)

returning *;
