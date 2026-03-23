-- name: UpdateBlog :one
update blogs set content = ? where id = ?

returning *;
