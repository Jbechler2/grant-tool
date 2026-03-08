-- name: GetUserByID :one
SELECT id, email, password_hash, role, created_at, updated_at, last_login_at, deleted_at
FROM users
WHERE id = $1
AND deleted_at is NULL;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, role, created_at, updated_at, last_login_at, deleted_at
FROM users
WHERE email = $1
AND deleted_at is NULL;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING id, email, password_hash, role, created_at, updated_at, last_login_at, deleted_at;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = NOW()
WHERE id=$1;