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

-- name: GetFeed :one
SELECT * from feeds
WHERE url = $1;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $1
WHERE feeds.id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
