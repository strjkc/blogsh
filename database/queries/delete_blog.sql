-- name: Delete :one
delete from blogs where id = ?

returning *;
