-- name: CreateClient :one
INSERT INTO clients (name, grant_writer_id)
VALUES ($1, $2)
RETURNING id, grant_writer_id, name, contact_name, contact_phone, contact_email, notes, created_at, updated_at, deleted_at;

-- name: GetClientByID :one
SELECT id, grant_writer_id, name, contact_name, contact_phone, contact_email, notes, created_at, updated_at, deleted_at
FROM clients
WHERE grant_writer_id = $1
AND id = $2
AND deleted_at IS NULL;

-- name: GetAllClientsByGrantWriter :many
SELECT id, grant_writer_id, name, contact_name, contact_phone, contact_email, notes, created_at, updated_at, deleted_at
FROM clients
WHERE grant_writer_id = $1
AND deleted_at IS NULL;

-- name: UpdateClient :one
UPDATE clients
SET 
  name = $2,
  contact_name = $3,
  contact_phone = $4,
  contact_email = $5,
  notes = $6,
  updated_at = NOW()
where id = $7
AND grant_writer_id = $1
AND deleted_at IS NULL
RETURNING id, grant_writer_id, name, contact_name, contact_phone, contact_email, notes, created_at, updated_at, deleted_at;

-- name: DeleteClient :exec
UPDATE clients
SET deleted_at = NOW()
WHERE grant_writer_id = $1
AND id=$2
AND deleted_at IS NULL;