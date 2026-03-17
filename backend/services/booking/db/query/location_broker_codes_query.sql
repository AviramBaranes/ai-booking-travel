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