-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5, 
    $6
)
RETURNING *;


-- name: GetFeeds :many
SELECT f.*, u.name as author_name FROM feeds f
left join users u
on u.id = f.user_id;

-- name: DeleteFeeds :exec
DELETE FROM feeds;
