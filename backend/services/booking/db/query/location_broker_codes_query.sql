-- name: GetLocationBrokerCode :one
SELECT
    id,
    location_id,
    broker,
    broker_location_id,
    enabled,
    created_at,
    updated_at
FROM
    location_broker_codes
WHERE
    (
        broker = sqlc.arg (broker)
        AND broker_location_id = sqlc.arg (broker_location_id)
    )
    OR (
        location_id = sqlc.arg (location_id)
        AND broker = sqlc.arg (broker)
    )
LIMIT
    1;

-- name: InsertLocationBrokerCode :one
INSERT INTO
    location_broker_codes (location_id, broker, broker_location_id)
VALUES
    (
        sqlc.arg (location_id),
        sqlc.arg (broker),
        sqlc.arg (broker_location_id)
    ) RETURNING id,
    location_id,
    broker,
    broker_location_id,
    enabled,
    created_at,
    updated_at;

-- name: EnableLocationBrokerCode :exec
UPDATE location_broker_codes
SET
    enabled = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg (id);

-- name: DisableLocationBrokerCode :exec
UPDATE location_broker_codes
SET
    enabled = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg (id);

-- name: GetAllLocationBrokerCodesByLocationIDs :many
SELECT
    lbc.*,
    l.country_code AS location_country_code
FROM
    location_broker_codes lbc
    JOIN locations l ON l.id = lbc.location_id
WHERE
    lbc.location_id = ANY (sqlc.arg ('location_ids')::bigint[]);

-- name: ListLocationBrokerCodesWithLocation :many
SELECT
    lbc.id,
    lbc.location_id,
    lbc.broker,
    lbc.broker_location_id,
    lbc.enabled,
    lbc.created_at,
    lbc.updated_at,
    l.country AS location_country,
    l.country_code AS location_country_code,
    l.city AS location_city,
    l.name AS location_name,
    l.iata AS location_iata
FROM
    location_broker_codes lbc
    JOIN locations l ON l.id = lbc.location_id
WHERE
    (sqlc.narg('country_code')::text IS NULL OR l.country_code ILIKE '%' || sqlc.narg('country_code')::text || '%')
    AND (sqlc.narg('broker')::text IS NULL OR lbc.broker::text ILIKE '%' || sqlc.narg('broker')::text || '%')
    AND (sqlc.narg('name')::text IS NULL OR l.name ILIKE '%' || sqlc.narg('name')::text || '%')
    AND (sqlc.narg('iata')::text IS NULL OR l.iata ILIKE '%' || sqlc.narg('iata')::text || '%')
    AND (sqlc.narg('enabled')::boolean IS NULL OR lbc.enabled = sqlc.narg('enabled')::boolean)
ORDER BY
    l.country_code, l.name, lbc.broker
LIMIT $1
OFFSET $2;

-- name: CountLocationBrokerCodesWithLocation :one
SELECT COUNT(*) AS total
FROM
    location_broker_codes lbc
    JOIN locations l ON l.id = lbc.location_id
WHERE
    (sqlc.narg('country_code')::text IS NULL OR l.country_code ILIKE '%' || sqlc.narg('country_code')::text || '%')
    AND (sqlc.narg('broker')::text IS NULL OR lbc.broker::text ILIKE '%' || sqlc.narg('broker')::text || '%')
    AND (sqlc.narg('name')::text IS NULL OR l.name ILIKE '%' || sqlc.narg('name')::text || '%')
    AND (sqlc.narg('iata')::text IS NULL OR l.iata ILIKE '%' || sqlc.narg('iata')::text || '%')
    AND (sqlc.narg('enabled')::boolean IS NULL OR lbc.enabled = sqlc.narg('enabled')::boolean);

-- name: DeleteLocationBrokerCode :one
DELETE FROM location_broker_codes
WHERE id = sqlc.arg(id)
RETURNING location_id;

-- name: CountLocationBrokerCodesByLocationID :one
SELECT COUNT(*) FROM location_broker_codes
WHERE location_id = sqlc.arg(location_id);