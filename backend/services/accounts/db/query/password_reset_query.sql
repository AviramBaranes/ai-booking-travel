-- name: InsertPasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, token_hash, expires_at, created_at;

-- name: GetPasswordResetTokenByHash :one
SELECT id, user_id, token_hash, expires_at, created_at
FROM password_reset_tokens
WHERE token_hash = $1;

-- name: DeletePasswordResetTokenByID :exec
DELETE FROM password_reset_tokens
WHERE id = $1;