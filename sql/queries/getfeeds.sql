-- name: GetFeeds :many
SELECT feeds.name AS feedName,
    feeds.url AS feedURL,
    users.name AS userName
FROM feeds
JOIN users
    ON feeds.user_id = users.id
ORDER BY feeds.created_at;
