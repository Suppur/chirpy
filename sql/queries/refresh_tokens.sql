-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
)
RETURNING token;

-- name: RevokeToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: GetUserByToken :one
SELECT user_id FROM refresh_tokens
WHERE token = $1;

-- name: GetToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;