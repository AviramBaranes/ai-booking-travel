-- name: ListOrganizations :many
SELECT
    o.id,
    o.name,
    o.is_organic,
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
    (sqlc.narg(name)::VARCHAR IS NULL       OR o.name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(is_organic)::BOOLEAN IS NULL OR o.is_organic = sqlc.narg(is_organic)::BOOLEAN)
GROUP BY o.id
ORDER BY o.name
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountOrganizations :one
SELECT COUNT(*)::BIGINT AS total
FROM organizations o
WHERE
    (sqlc.narg(name)::VARCHAR IS NULL       OR o.name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(is_organic)::BOOLEAN IS NULL OR o.is_organic = sqlc.narg(is_organic)::BOOLEAN);

-- name: CreateOrganization :one
INSERT INTO organizations (name, is_organic, phone, address, obligo, created_at, updated_at)
VALUES (
    sqlc.arg(name),
    sqlc.arg(is_organic),
    sqlc.narg(phone),
    sqlc.narg(address),
    sqlc.narg(obligo),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
RETURNING id, name, is_organic, phone, address, obligo, created_at, updated_at;

-- name: UpdateOrganization :one
UPDATE organizations
SET
    name       = COALESCE(sqlc.narg(name),       name),
    is_organic = COALESCE(sqlc.narg(is_organic), is_organic),
    phone      = COALESCE(sqlc.narg(phone),      phone),
    address    = COALESCE(sqlc.narg(address),    address),
    obligo     = COALESCE(sqlc.narg(obligo),     obligo),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, name, is_organic, phone, address, obligo, created_at, updated_at;

-- name: ListOrganicOrganizations :many
SELECT
    id,name
FROM organizations
WHERE is_organic = TRUE
ORDER BY name;
