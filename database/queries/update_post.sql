-- name: UpdatePost :one
update posts set content = ?, title = ?, tags = ?, category = ?, updatedat = ?  where id = ?

returning *;
