-- name: GetPostsFilter :many
select * from posts 
WHERE title LIKE '%' || ? || '%'
   OR content LIKE '%' || ? || '%'
   OR category LIKE '%' || ? || '%';
