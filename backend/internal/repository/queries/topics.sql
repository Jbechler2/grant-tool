-- name: CreateTopic :one
INSERT INTO topics (grant_writer_id, label)
VALUES ($1, $2)
RETURNING id, grant_writer_id, label;

-- name: GetAllTopics :many
SELECT id, label
FROM topics
where grant_writer_id = $1;

-- name: GetAllTopicsByGrant :many
SELECT id, label
FROM grants_topics gt
JOIN topics t ON t.id=gt.topic_id
WHERE t.grant_writer_id = $1
AND gt.grant_id = $2;

-- name: GetAllTopicsByClient :many
SELECT id, label
FROM clients_topics ct
JOIN topics t ON t.id=ct.topic_id
WHERE t.grant_writer_id = $1
AND ct.client_id = $2;

-- name: AddTopicToGrant :execrows
INSERT INTO grants_topics (topic_id, grant_id)
SELECT $2, g.id
FROM grants g
WHERE g.id = $1
AND g.grant_writer_id = $3;

-- name: AddTopicToClient :execrows
INSERT INTO clients_topics (topic_id, client_id)
SELECT $2, c.id
FROM clients c
WHERE c.id = $1
AND c.grant_writer_id = $3;

-- name: UpdateTopic :one
UPDATE topics
SET label = $3
WHERE id = $1
AND grant_writer_id = $2
RETURNING id, label;

-- name: DeleteGrantTopic :exec
DELETE FROM grants_topics gt
USING topics t
WHERE gt.topic_id = t.id
AND gt.topic_id = $1
AND gt.grant_id = $2
AND t.grant_writer_id = $3;

-- name: DeleteClientTopic :exec
DELETE FROM clients_topics ct
USING topics t
WHERE ct.topic_id = t.id
AND ct.topic_id = $1
AND ct.client_id = $2
AND t.grant_writer_id = $3;

-- name: DeleteTopic :exec
DELETE FROM topics
where id = $1
and grant_writer_id = $2;