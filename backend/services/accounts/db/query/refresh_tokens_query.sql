-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE jti = $1;

-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens (jti, user_id, admin_ref_id, expires_at)
VALUES ($1, $2, $3, $4);

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE jti = $1;

-- name: DeleteRefreshTokensByUserId :exec
DELETE FROM refresh_tokens WHERE user_id = $1;