-- name: CreateAgent :one
INSERT INTO users (role, email, phone_number, password_hash, office_id, created_at, updated_at)
VALUES ('agent', sqlc.arg(email), sqlc.arg(phone_number), sqlc.arg(password_hash), sqlc.arg(office_id), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, email, phone_number, office_id, last_login, created_at, updated_at;

-- name: CreateAdmin :one
INSERT INTO users (role, email, password_hash, created_at, updated_at)
VALUES ('admin', sqlc.arg(email), sqlc.arg(password_hash), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, email, office_id, last_login, created_at, updated_at;

-- name: CreateCustomer :one
INSERT INTO users (role, email, phone_number, otp, password_hash, created_at, updated_at)
VALUES (
  'customer',
  sqlc.arg(email),
  sqlc.arg(phone_number),
  sqlc.narg(otp)::varchar,
  sqlc.arg(password_hash),
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
)
RETURNING id, role, email, phone_number, otp, office_id, last_login, created_at, updated_at;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CheckUserExists :one
SELECT id FROM users
WHERE email = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone_number = $1;

-- name: ListAgents :many
SELECT u.id, u.role, u.email, u.phone_number, u.office_id, u.last_login, u.created_at, u.updated_at,
       o.name AS office_name,
       org.name AS organization_name
FROM users u
LEFT JOIN offices       o   ON o.id   = u.office_id
LEFT JOIN organizations org ON org.id = o.organization_id
WHERE u.role = 'agent'
  AND (sqlc.narg(search)::text IS NULL OR u.email ILIKE '%' || sqlc.narg(search)::text || '%' OR u.phone_number ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(office_id)::int IS NULL OR u.office_id = sqlc.narg(office_id)::int)
  AND (sqlc.narg(organization_id)::int IS NULL OR o.organization_id = sqlc.narg(organization_id)::int)
ORDER BY u.created_at DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);

-- name: CountAgents :one
SELECT COUNT(*)
FROM users
WHERE role = 'agent'
  AND (sqlc.narg(search)::text IS NULL OR email ILIKE '%' || sqlc.narg(search)::text || '%' OR phone_number ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(office_id)::int IS NULL OR office_id = sqlc.narg(office_id)::int)
  AND (sqlc.narg(organization_id)::int IS NULL OR office_id IN (
    SELECT id FROM offices WHERE organization_id = sqlc.narg(organization_id)::int
  ));

-- name: ListAdmins :many
SELECT id, role, email, office_id, last_login, created_at, updated_at
FROM users
WHERE role = 'admin';

-- name: UpdateUser :one
UPDATE users
SET
  email = COALESCE(sqlc.narg(email)::varchar, email),
  phone_number = COALESCE(sqlc.narg(phone_number)::varchar, phone_number),
  office_id = COALESCE(sqlc.narg(office_id)::int, office_id),
  password_hash = COALESCE(sqlc.narg(password_hash)::varchar, password_hash),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, role, email, phone_number, office_id, last_login, created_at, updated_at;

-- name: ListAdminsEmails :many
SELECT email
FROM users
WHERE role = 'admin';

-- name: SaveOTP :exec
UPDATE users
SET
  otp = $2
WHERE
  id = $1;

-- name: GetAgentsBillingContacts :many
SELECT
u.id as agent_id,
c.email, c.first_name, c.last_name,
org.id as organization_id, org.name as organization_name, 
office.id as office_id, office.name as office_name
FROM users as u
INNER JOIN offices as office ON office.id = u.office_id
INNER JOIN organizations as org ON org.id = office.organization_id
INNER JOIN contacts as c ON c.organization_id = org.id AND c.is_payment_responsible = TRUE
WHERE u.role = 'agent'
  AND u.id = ANY(sqlc.arg(users_ids)::int[]);
