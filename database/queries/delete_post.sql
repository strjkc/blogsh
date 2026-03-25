-- name: Delete :one
delete from posts where id = ?

returning *;
