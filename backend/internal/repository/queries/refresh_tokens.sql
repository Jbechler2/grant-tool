-- name: CreateToken :one
INSERT INTO refresh_tokens (grant_writer_id, token, user_agent, ip_address)
VALUES ($1, $2, $3, $4)
RETURNING id, grant_writer_id, token, user_agent, ip_address, created_at, expires_at;

-- name: GetRefreshTokenByTokenValue :one
SELECT id, grant_writer_id, token, user_agent, ip_address, created_at, expires_at
FROM refresh_tokens
WHERE token = $1;

-- name: CountValidTokens :one
SELECT count(*)
FROM refresh_tokens
WHERE grant_writer_id = $1
AND expires_at > NOW();

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE grant_writer_id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();