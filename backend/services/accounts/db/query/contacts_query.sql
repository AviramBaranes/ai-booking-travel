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
    (sqlc.narg(name)::VARCHAR IS NULL OR c.first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR c.last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(office_id)::INTEGER IS NULL       OR c.office_id = sqlc.narg(office_id)::INTEGER)
    AND (sqlc.narg(organization_id)::INTEGER IS NULL OR c.organization_id = sqlc.narg(organization_id)::INTEGER)
ORDER BY c.last_name, c.first_name
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountContacts :one
SELECT COUNT(*)::BIGINT AS total
FROM contacts c
LEFT JOIN offices o ON o.id = c.office_id
WHERE
    (sqlc.narg(name)::VARCHAR IS NULL OR c.first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR c.last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(office_id)::INTEGER IS NULL       OR c.office_id = sqlc.narg(office_id)::INTEGER)
    AND (sqlc.narg(organization_id)::INTEGER IS NULL OR c.organization_id = sqlc.narg(organization_id)::INTEGER);

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
    first_name      = COALESCE(sqlc.narg(first_name),      first_name),
    last_name       = COALESCE(sqlc.narg(last_name),       last_name),
    role            = COALESCE(sqlc.narg(role),            role),
    cellphone       = COALESCE(sqlc.narg(cellphone),       cellphone),
    email           = COALESCE(sqlc.narg(email),           email),
    office_id              = COALESCE(sqlc.narg(office_id),              office_id),
    organization_id        = COALESCE(sqlc.narg(organization_id),        organization_id),
    is_payment_responsible = COALESCE(sqlc.narg(is_payment_responsible), is_payment_responsible),
    updated_at             = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, first_name, last_name, role, cellphone, email, office_id, organization_id, is_payment_responsible, created_at, updated_at;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = sqlc.arg(id);
