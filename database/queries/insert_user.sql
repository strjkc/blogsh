-- name: InsertUser :one
insert into users(username, password) values(?, ?)

RETURNING *;
