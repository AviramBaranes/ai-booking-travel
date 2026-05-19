-- name: ListOrganizations :many
SELECT
    o.id,
    o.name,
    o.is_organic,
    o.icount_client_id,
    o.phone,
    o.address,
    o.obligo,
    o.created_at,
    o.updated_at,
    COUNT(DISTINCT of.id)::BIGINT          AS office_count,
    COUNT(DISTINCT c.id)::BIGINT           AS contact_count,
    COUNT(DISTINCT u.id)::BIGINT           AS agent_count
FROM organizations o
LEFT JOIN offices of ON of.organization_id = o.id
LEFT JOIN contacts c ON (c.organization_id = o.id OR c.office_id = of.id)
LEFT JOIN users u ON (u.office_id = of.id AND u.role = 'agent')
WHERE
    (sqlc.narg(name)::TEXT IS NULL       OR o.name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(is_organic)::BOOLEAN IS NULL OR o.is_organic = sqlc.narg(is_organic)::BOOLEAN)
GROUP BY o.id
ORDER BY o.name
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountOrganizations :one
SELECT COUNT(*)::BIGINT AS total
FROM organizations o
WHERE
    (sqlc.narg(name)::TEXT IS NULL       OR o.name ILIKE '%' || sqlc.narg(name)::TEXT || '%')
    AND (sqlc.narg(is_organic)::BOOLEAN IS NULL OR o.is_organic = sqlc.narg(is_organic)::BOOLEAN);

-- name: CreateOrganization :one
INSERT INTO organizations (name, is_organic, icount_client_id, phone, address, obligo, created_at, updated_at)
VALUES (
    sqlc.arg(name),
    sqlc.arg(is_organic),
    sqlc.narg(icount_client_id),
    sqlc.narg(phone),
    sqlc.narg(address),
    sqlc.narg(obligo),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING id, name, is_organic, icount_client_id, phone, address, obligo, created_at, updated_at;

-- name: UpdateOrganization :one
UPDATE organizations
SET
    name             = sqlc.arg(name),
    is_organic       = sqlc.arg(is_organic),
    icount_client_id = sqlc.narg(icount_client_id),
    phone            = sqlc.narg(phone),
    address          = sqlc.narg(address),
    obligo           = sqlc.narg(obligo),
    updated_at       = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, name, is_organic, icount_client_id, phone, address, obligo, created_at, updated_at;

-- name: ListOrganicOrganizations :many
SELECT
    id,name
FROM organizations
WHERE is_organic = TRUE
ORDER BY name;

-- name: GetOrganizationIcountClientID :one
SELECT icount_client_id
FROM organizations
WHERE id = sqlc.arg(id)::BIGINT;

-- name: GetOrganizationBillingState :one
SELECT is_organic, icount_client_id
FROM organizations
WHERE id = sqlc.arg(id)::BIGINT;
