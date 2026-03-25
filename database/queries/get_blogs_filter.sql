-- name: GetBlogsFilter :many
select * from blogs 
WHERE title LIKE '%' || ? || '%'
   OR content LIKE '%' || ? || '%'
   OR category LIKE '%' || ? || '%';
