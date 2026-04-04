-- name: RegisterAgent :one
INSERT INTO users (role, email, phone_number, password_hash, office_id, created_at, updated_at)
VALUES ('agent', sqlc.arg(email), sqlc.arg(phone_number), sqlc.arg(password_hash), sqlc.arg(office_id), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, email, phone_number, office_id, last_login, created_at, updated_at;

-- name: RegisterAdmin :one
INSERT INTO users (role, email, password_hash, created_at, updated_at)
VALUES ('admin', sqlc.arg(email), sqlc.arg(password_hash), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, email, office_id, last_login, created_at, updated_at;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CheckUserExists :one
SELECT id FROM users
WHERE email = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateUser :one
UPDATE users
SET
  office_id    = COALESCE(sqlc.narg(office_id),    office_id),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  updated_at   = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, role, email, office_id, phone_number, last_login, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
