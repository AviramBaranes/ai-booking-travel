-- name: ListContacts :many
SELECT
    c.id,
    c.first_name,
    c.last_name,
    c.role,
    c.cellphone,
    c.email,
    c.office_id,
    c.organization_id,
    c.is_payment_responsible,
    c.created_at,
    c.updated_at,
    o.name  AS office_name,
    org.name AS organization_name
FROM contacts c
LEFT JOIN offices       o   ON o.id   = c.office_id
LEFT JOIN organizations org ON org.id = c.organization_id
WHERE
    (sqlc.narg(name)::TEXT IS NULL OR c.first_name ILIKE '%' || sqlc.narg(name)::TEXT || '%' OR c.last_name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(office_id)::BIGINT IS NULL       OR c.office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR c.organization_id = sqlc.narg(organization_id)::BIGINT)
ORDER BY c.last_name, c.first_name
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountContacts :one
SELECT COUNT(*)::BIGINT AS total
FROM contacts c
LEFT JOIN offices o ON o.id = c.office_id
WHERE
    (sqlc.narg(name)::TEXT IS NULL OR c.first_name ILIKE '%' || sqlc.narg(name)::TEXT || '%' OR c.last_name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(office_id)::BIGINT IS NULL       OR c.office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR c.organization_id = sqlc.narg(organization_id)::BIGINT);

-- name: CreateContact :one
INSERT INTO contacts (first_name, last_name, role, cellphone, email, office_id, organization_id, is_payment_responsible, created_at, updated_at)
VALUES (
    sqlc.arg(first_name),
    sqlc.arg(last_name),
    sqlc.arg(role),
    sqlc.arg(cellphone),
    sqlc.arg(email),
    sqlc.narg(office_id),
    sqlc.narg(organization_id),
    sqlc.arg(is_payment_responsible),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING id, first_name, last_name, role, cellphone, email, office_id, organization_id, is_payment_responsible, created_at, updated_at;

-- name: UpdateContact :one
UPDATE contacts
SET
    first_name             = sqlc.arg(first_name),
    last_name              = sqlc.arg(last_name),
    role                   = sqlc.arg(role),
    cellphone              = sqlc.arg(cellphone),
    email                  = sqlc.arg(email),
    office_id              = sqlc.narg(office_id),
    organization_id        = sqlc.narg(organization_id),
    is_payment_responsible = sqlc.arg(is_payment_responsible),
    updated_at             = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, first_name, last_name, role, cellphone, email, office_id, organization_id, is_payment_responsible, created_at, updated_at;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = sqlc.arg(id);
