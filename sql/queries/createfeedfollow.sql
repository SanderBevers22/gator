-- name: CreateFeedFollow :one
WITH inserted_feed_follow as (
    INSERT INTO feed_follows (id,created_at,updated_at,user_id,feed_id)
    VALUES ($1,$2,$3,$4,$5)
    returning *
)
SELECT
    inserted_feed_follow.*,feeds.name AS feedname, users.name AS username
FROM inserted_feed_follow
INNER JOIN feeds
    ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users
    ON inserted_feed_follow.user_id = users.id;
