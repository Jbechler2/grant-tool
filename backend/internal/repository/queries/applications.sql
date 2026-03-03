-- name: CreateApplication :one
INSERT INTO applications (grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at;

-- name: GetApplicationByID :one
SELECT id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at
FROM applications
WHERE deleted_at IS NULL
AND grant_writer_id = $1
and id = $2;

-- name: GetAllApplicationsByUserID :many
SELECT id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at
FROM applications
WHERE deleted_at IS NULL
AND grant_writer_id = $1;

-- name: GetAllApplicationsByClientID :many
SELECT id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at
FROM applications
WHERE deleted_at IS NULL
AND grant_writer_id = $1
AND client_id = $2;

-- name: UpdateApplication :one
UPDATE applications
SET
  title = $3,
  status = $4,
  is_exclusive = $5,
  notes = $6,
  updated_at = NOW()
WHERE deleted_at IS NULL
AND id = $1
AND grant_writer_id = $2
RETURNING id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at;

-- name: DeleteApplication :exec
UPDATE applications
SET deleted_at = NOW()
WHERE deleted_at IS NULL
AND grant_writer_id = $1
AND id = $2;

-- name: PublishApplication :one
UPDATE applications
SET 
  published_at = now()
WHERE deleted_at IS NULL
AND published_at IS NULL
AND grant_writer_id = $1
AND id = $2
RETURNING id, grant_writer_id, grant_id, client_id, title, status, is_exclusive, published_at, notes, created_at, updated_at, deleted_at;