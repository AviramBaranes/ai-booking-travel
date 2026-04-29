-- name: GetOutboxByTopic :many
SELECT id, topic, data, inserted_at
FROM outbox
WHERE topic = @topic
ORDER BY id DESC;
