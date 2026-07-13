-- name: CreateAgent :one
INSERT INTO users (role, first_name, last_name, email, phone_number, password_hash, office_id, created_at, updated_at)
VALUES ('agent', sqlc.arg(first_name), sqlc.arg(last_name), sqlc.arg(email), sqlc.arg(phone_number), sqlc.arg(password_hash), sqlc.arg(office_id), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, first_name, last_name, email, phone_number, office_id, last_login, created_at, updated_at;

-- name: CreateStaffUser :one
INSERT INTO users (role, first_name, last_name, email, password_hash, created_at, updated_at)
VALUES (sqlc.arg(role)::user_role, sqlc.arg(first_name), sqlc.arg(last_name), sqlc.arg(email), sqlc.arg(password_hash), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, role, first_name, last_name, email, office_id, last_login, created_at, updated_at;

-- name: CreateCustomer :one
INSERT INTO users (role, first_name, last_name, email, phone_number, otp, password_hash, created_at, updated_at)
VALUES (
  'customer',
  sqlc.arg(first_name),
  sqlc.arg(last_name),
  sqlc.arg(email),
  sqlc.arg(phone_number),
  sqlc.narg(otp)::text,
  sqlc.arg(password_hash),
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
)
RETURNING id, role, first_name, last_name, email, phone_number, otp, office_id, last_login, created_at, updated_at;

-- name: GetUserById :one
SELECT users.*, organization.id AS organization_id, organization.is_organic
FROM users
LEFT JOIN offices ON offices.id = users.office_id
LEFT JOIN organizations AS organization ON organization.id = offices.organization_id
WHERE users.id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CheckUserExists :one
SELECT id FROM users
WHERE email = $1;

-- name: GetUserByEmail :one
SELECT users.*, organization.id AS organization_id, organization.is_organic
FROM users
LEFT JOIN offices ON offices.id = users.office_id
LEFT JOIN organizations AS organization ON organization.id = offices.organization_id
WHERE email = $1;

-- name: GetUserByPhone :one
SELECT *
FROM users
WHERE phone_number = $1;

-- name: ListAgents :many
SELECT u.id, u.role, u.first_name, u.last_name, u.email, u.phone_number, u.office_id, u.last_login, u.created_at, u.updated_at,
       o.name AS office_name,
       org.name AS organization_name
FROM users u
LEFT JOIN offices       o   ON o.id   = u.office_id
LEFT JOIN organizations org ON org.id = o.organization_id
WHERE u.role = 'agent'
  AND (sqlc.narg(search)::text IS NULL 
    OR u.email ILIKE '%' || sqlc.narg(search)::text || '%' 
    OR u.phone_number ILIKE '%' || sqlc.narg(search)::text || '%'
    OR u.first_name ILIKE '%' || sqlc.narg(search)::text || '%'
    OR u.last_name ILIKE '%' || sqlc.narg(search)::text || '%'
    )
  AND (sqlc.narg(office_id)::bigint IS NULL OR u.office_id = sqlc.narg(office_id)::bigint)
  AND (sqlc.narg(organization_id)::bigint IS NULL OR o.organization_id = sqlc.narg(organization_id)::bigint)
ORDER BY u.created_at DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);

-- name: CountAgents :one
SELECT COUNT(*)
FROM users
WHERE role = 'agent'
  AND (sqlc.narg(search)::text IS NULL OR email ILIKE '%' || sqlc.narg(search)::text || '%' OR phone_number ILIKE '%' || sqlc.narg(search)::text || '%')
  AND (sqlc.narg(office_id)::bigint IS NULL OR office_id = sqlc.narg(office_id)::bigint)
  AND (sqlc.narg(organization_id)::bigint IS NULL OR office_id IN (
    SELECT id FROM offices WHERE organization_id = sqlc.narg(organization_id)::bigint
  ));

-- name: ListStaffByRole :many
SELECT id, role, first_name, last_name, email, office_id, last_login, created_at, updated_at
FROM users
WHERE role = sqlc.arg(role)::user_role;

-- name: ListAdminsEmails :many
SELECT email
FROM users
WHERE role = 'admin';

-- name: GetUserNamesByIDs :many
SELECT id, first_name, last_name
FROM users
WHERE id = ANY(sqlc.arg(ids)::BIGINT[]);

-- name: UpdateUser :one
UPDATE users
SET
  first_name = COALESCE(sqlc.narg(first_name)::text, first_name),
  last_name = COALESCE(sqlc.narg(last_name)::text, last_name),
  email = COALESCE(sqlc.narg(email)::text, email),
  phone_number = COALESCE(sqlc.narg(phone_number)::text, phone_number),
  office_id = COALESCE(sqlc.narg(office_id)::bigint, office_id),
  password_hash = COALESCE(sqlc.narg(password_hash)::text, password_hash),
  updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, role, first_name, last_name, email, phone_number, office_id, last_login, created_at, updated_at;

-- name: SaveOTP :exec
UPDATE users
SET
  otp = $2
WHERE
  id = $1;

-- name: GetAgentsBillingContacts :many
SELECT
u.id as agent_id, u.first_name as agent_first_name, u.last_name as agent_last_name,
c.id as contact_id, c.email, c.first_name as contact_first_name, c.last_name as contact_last_name,
org.id as organization_id, org.name as organization_name, org.is_organic,
office.id as office_id, office.name as office_name
FROM users as u
INNER JOIN offices as office ON office.id = u.office_id
INNER JOIN organizations as org ON org.id = office.organization_id
INNER JOIN contacts as c ON c.is_payment_responsible = TRUE AND (
    (c.organization_id = org.id AND org.is_organic = TRUE) OR 
    (c.office_id = office.id AND org.is_organic = FALSE)
)
WHERE u.role = 'agent'
  AND u.id = ANY(sqlc.arg(users_ids)::bigint[]);

-- name: GetUserCredit :one
SELECT (
  CASE WHEN org.is_organic = TRUE
    THEN COALESCE(org.obligo, 0)
    ELSE COALESCE(offices.obligo, 0)
  END
)::INTEGER AS obligo,
(
  CASE WHEN org.is_organic = TRUE
    THEN org.balance_due
    ELSE offices.balance_due
  END
)::NUMERIC(10, 2) AS balance_due
FROM users AS u
INNER JOIN offices ON u.office_id = offices.id
INNER JOIN organizations AS org ON org.id = offices.organization_id
WHERE u.id = $1;

-- name: GetUserGrossMarkup :one
SELECT 
  COALESCE(
    NULLIF(o.gross_markup, 0),
    org.gross_markup,
    0
  )::NUMERIC(5, 2) AS gross_markup
FROM users AS u
LEFT JOIN offices AS o ON u.office_id = o.id
LEFT JOIN organizations AS org ON org.id = o.organization_id
WHERE u.id = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetCustomerByPhoneAndEmail :one
SELECT
  id,
  role,
  first_name,
  last_name,
  email,
  phone_number,
  office_id,
  last_login,
  created_at,
  updated_at
FROM users
WHERE role = 'customer'
  AND phone_number = $1
  AND email = $2;

-- name: GetCustomerByID :one
SELECT * FROM users
WHERE role = 'customer'
  AND id = $1;

-- name: ListCustomers :many
SELECT u.id, u.role, u.first_name, u.last_name, u.email, u.phone_number, u.last_login, u.created_at, u.updated_at
FROM users u
WHERE u.role = 'customer'
  AND (sqlc.narg(search)::text IS NULL 
    OR u.email ILIKE '%' || sqlc.narg(search)::text || '%' 
    OR u.phone_number ILIKE '%' || sqlc.narg(search)::text || '%'
    OR u.first_name ILIKE '%' || sqlc.narg(search)::text || '%'
    OR u.last_name ILIKE '%' || sqlc.narg(search)::text || '%'
    )
ORDER BY u.created_at DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);

-- name: CountCustomers :one
SELECT COUNT(*)
FROM users
WHERE role = 'customer'
  AND (
        sqlc.narg(search)::text IS NULL 
        OR email ILIKE '%' || sqlc.narg(search)::text || '%' 
        OR phone_number ILIKE '%' || sqlc.narg(search)::text || '%'
      );