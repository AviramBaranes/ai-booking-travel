-- name: ListOffices :many
SELECT
    o.id,
    o.name,
    o.organization_id,
    org.name AS organization_name,
    o.icount_client_id,
    o.phone,
    o.address,
    o.created_at,
    o.updated_at,
    COUNT(DISTINCT c.id)::BIGINT  AS contact_count,
    COUNT(DISTINCT u.id)::BIGINT  AS agent_count
FROM offices o
JOIN organizations org ON org.id = o.organization_id
LEFT JOIN contacts c ON c.office_id = o.id
LEFT JOIN users u ON (u.office_id = o.id AND u.role = 'agent')
WHERE
    (sqlc.narg(name)::TEXT IS NULL            OR o.name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR o.organization_id = sqlc.narg(organization_id)::BIGINT)
GROUP BY o.id, org.name
ORDER BY o.name
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountOffices :one
SELECT COUNT(*)::BIGINT AS total
FROM offices o
WHERE
    (sqlc.narg(name)::TEXT IS NULL            OR o.name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR o.organization_id = sqlc.narg(organization_id)::BIGINT);

-- name: CreateOffice :one
INSERT INTO offices (name, organization_id, icount_client_id, phone, address, created_at, updated_at)
VALUES (
    sqlc.arg(name),
    sqlc.arg(organization_id),
    sqlc.narg(icount_client_id),
    sqlc.narg(phone),
    sqlc.narg(address),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING id, name, organization_id, icount_client_id, phone, address, created_at, updated_at;

-- name: UpdateOffice :one
UPDATE offices
SET
    name             = sqlc.arg(name),
    organization_id  = sqlc.arg(organization_id),
    icount_client_id = sqlc.narg(icount_client_id),
    phone            = sqlc.narg(phone),
    address          = sqlc.narg(address),
    updated_at       = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, name, organization_id, icount_client_id, phone, address, created_at, updated_at;

-- name: ListInorganicOffices :many
SELECT
    o.id, o.name
FROM offices as o
JOIN organizations org ON org.id = o.organization_id
WHERE org.is_organic = FALSE
ORDER BY o.name;

-- name: GetOfficeIcountClientID :one
SELECT icount_client_id
FROM offices
WHERE id = sqlc.arg(id)::BIGINT;

-- name: GetOfficeBillingState :one
SELECT o.icount_client_id, o.organization_id, org.is_organic
FROM offices o
JOIN organizations org ON org.id = o.organization_id
WHERE o.id = sqlc.arg(id)::BIGINT;