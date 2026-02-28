-- name: CreateGrant :one
INSERT INTO grants (grant_writer_id, title, funder_name, visibility, funder_website, description, award_amount_min, award_amount_max, eligibility_notes, estimated_application_hours)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, grant_writer_id, title, funder_name, funder_website, description, award_amount_min, award_amount_max, eligibility_notes, estimated_application_hours, visibility, created_at, updated_at, deleted_at;

-- name: GetAllGrants :many
SELECT id, grant_writer_id, title, funder_name, funder_website, description, award_amount_min, award_amount_max, eligibility_notes, estimated_application_hours, visibility, created_at, updated_at, deleted_at
FROM grants
WHERE grant_writer_id = $1
AND deleted_at IS NULL;

-- name: GetGrantByID :one
SELECT id, grant_writer_id, title, funder_name, funder_website, description, award_amount_min, award_amount_max, eligibility_notes, estimated_application_hours, visibility, created_at, updated_at, deleted_at
FROM grants
WHERE grant_writer_id = $1
AND id = $2
AND deleted_at IS NULL;

-- name: UpdateGrant :one
UPDATE grants
set
  title = $3,
  funder_name = $4,
  visibility = $5,
  funder_website = $6,
  description = $7,
  award_amount_min = $8,
  award_amount_max = $9,
  eligibility_notes = $10,
  estimated_application_hours = $11,
  updated_at = NOW()
WHERE grant_writer_id = $1
AND  id = $2
AND deleted_at IS NULL
RETURNING id, grant_writer_id, title, funder_name, funder_website, description, award_amount_min, award_amount_max, eligibility_notes, estimated_application_hours, visibility, created_at, updated_at, deleted_at;

-- name: DeleteGrant :exec
UPDATE grants
SET deleted_at = NOW()
WHERE grant_writer_id = $1
AND id = $2
AND deleted_at IS NULL;

-- name: CreateDeadline :one
INSERT INTO grant_deadlines (grant_id, label, date, description)
VALUES ($1, $2, $3, $4)
RETURNING id, grant_id, label, date, description, created_at;


-- name: GetDeadlinesByGrantID :many
SELECT gd.id, grant_id, label, date, gd.description, gd.created_at
FROM grant_deadlines gd
JOIN grants g ON g.id = gd.grant_id
WHERE gd.grant_id = $2
AND g.grant_writer_id = $1
AND g.deleted_at IS NULL;

-- name: DeleteDeadline :exec
DELETE FROM grant_deadlines gd
USING grants g
WHERE g.id= gd.grant_id
AND gd.grant_id = $2
AND g.grant_writer_id = $1
AND gd.id = $3
AND g.deleted_at IS NULL;
