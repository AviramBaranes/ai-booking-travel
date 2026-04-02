-- name: RegisterAgent :one
INSERT INTO users (role, username, password_hash, office_code, agent_code, created_at, updated_at)
VALUES ('agent', $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, username, office_code, agent_code, last_login, created_at, updated_at;

-- name: RegisterAdmin :one
INSERT INTO users (role, username, password_hash, office_code, agent_code, created_at, updated_at)
VALUES ('admin', $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, username, office_code, agent_code, last_login, created_at, updated_at;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CheckUserExists :one
SELECT id FROM users
WHERE username = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2;

-- name: UpdateUser :one
UPDATE users
SET
  agent_code   = COALESCE(sqlc.narg(agent_code),   agent_code),
  office_code  = COALESCE(sqlc.narg(office_code),  office_code),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  updated_at   = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, role, username, office_code, agent_code, phone_number, last_login, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;