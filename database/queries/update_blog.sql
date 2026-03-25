-- name: UpdateBlog :one
update blogs set content = ?, title = ?, tags = ?, updatedat = ?  where id = ?

returning *;
