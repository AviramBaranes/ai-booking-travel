-- name: UpsertLocationByIATA :one
INSERT INTO locations (country, country_code, city, name, iata)
VALUES (
  NULLIF(sqlc.arg(country), '')::text,
  NULLIF(sqlc.arg(country_code), '')::text,
  NULLIF(sqlc.arg(city), '')::text,
  sqlc.arg(name)::text,
  NULLIF(upper(sqlc.arg(iata)), '')::char(3)
)
ON CONFLICT (iata) WHERE iata IS NOT NULL
DO UPDATE SET
  country      = EXCLUDED.country,
  country_code = EXCLUDED.country_code,
  city         = EXCLUDED.city,
  name         = EXCLUDED.name,
  updated_at   = now()
RETURNING id;

-- name: UpsertLocationByCountryCodeName :one
INSERT INTO locations (country, country_code, city, name, iata)
VALUES (
  NULLIF(sqlc.arg(country), '')::text,
  NULLIF(sqlc.arg(country_code), '')::text,
  NULLIF(sqlc.arg(city), '')::text,
  sqlc.arg(name)::text,
  NULL
)
ON CONFLICT (country_code, lower(name))
DO UPDATE SET
  country      = EXCLUDED.country,
  country_code = EXCLUDED.country_code,
  city         = EXCLUDED.city,
  name         = EXCLUDED.name,
  updated_at   = now()
RETURNING id;

-- name: InsertManyLocation :many
INSERT INTO locations (country, country_code, city, name, iata)
SELECT
  NULLIF(unnest(sqlc.arg(countries)::text[]), '')::text,
  NULLIF(unnest(sqlc.arg(country_codes)::text[]), '')::text,
  NULLIF(unnest(sqlc.arg(cities)::text[]), '')::text,
  unnest(sqlc.arg(names)::text[])::text,
  NULLIF(upper(unnest(sqlc.arg(iatas)::text[])), '')::char(3)
RETURNING id;

-- name: GetLocationById :one
SELECT *
FROM locations
WHERE id = $1;

-- name: InsertLocation :one
INSERT INTO locations (country, country_code, city, name, iata)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetLocationByBrokerLocationID :one
SELECT l.*
FROM locations l
JOIN location_broker_codes lbc ON lbc.location_id = l.id
WHERE lbc.broker_location_id = @broker_location_id
LIMIT 1;

-- name: GetLocationIDByBrokerCode :one
SELECT lbc.location_id
FROM location_broker_codes lbc
WHERE lbc.broker = sqlc.arg(broker)::broker
  AND lbc.broker_location_id = sqlc.arg(broker_location_id)
  AND lbc.enabled = TRUE
LIMIT 1;

-- name: DeleteLocationByID :exec
DELETE FROM locations
WHERE id = sqlc.arg(id);

-- name: SearchLocations :many
WITH matched_locations AS (
    SELECT
        l.id,
        l.country,
        l.country_code,
        l.city,
        l.name,
        l.iata,
        l.is_airport,
        l.created_at,
        l.updated_at,
        CASE
            -- IATA exact always wins
            WHEN upper(l.iata::text) = upper(sqlc.arg(search)::text) THEN 0

            -- Exact alias
            WHEN EXISTS (
                SELECT 1
                FROM location_aliases la
                WHERE la.location_id = l.id
                  AND lower(la.alias) = lower(sqlc.arg(search)::text)
            ) THEN 1

            -- Exact canonical fields
            WHEN lower(l.name) = lower(sqlc.arg(search)::text) THEN 2
            WHEN lower(coalesce(l.city, '')) = lower(sqlc.arg(search)::text) THEN 3
            WHEN lower(l.country) = lower(sqlc.arg(search)::text) THEN 4

            -- Prefix matches
            WHEN l.name ILIKE sqlc.arg(search)::text || '%' THEN 5
            WHEN coalesce(l.city, '') ILIKE sqlc.arg(search)::text || '%' THEN 6
            WHEN l.country ILIKE sqlc.arg(search)::text || '%' THEN 7

            WHEN EXISTS (
                SELECT 1
                FROM location_aliases la
                WHERE la.location_id = l.id
                  AND la.alias ILIKE sqlc.arg(search)::text || '%'
            ) THEN 8

            -- Contains fallback
            WHEN l.name ILIKE '%' || sqlc.arg(search)::text || '%' THEN 9
            WHEN coalesce(l.city, '') ILIKE '%' || sqlc.arg(search)::text || '%' THEN 10
            WHEN l.country ILIKE '%' || sqlc.arg(search)::text || '%' THEN 11
            WHEN l.iata ILIKE '%' || sqlc.arg(search)::text || '%' THEN 12

            WHEN EXISTS (
                SELECT 1
                FROM location_aliases la
                WHERE la.location_id = l.id
                  AND la.alias ILIKE '%' || sqlc.arg(search)::text || '%'
            ) THEN 13

            ELSE 14
        END AS search_rank
    FROM locations l
    WHERE EXISTS (
        SELECT 1
        FROM location_broker_codes lbc
        WHERE lbc.location_id = l.id
          AND lbc.enabled = TRUE
    )
    AND (
        l.name ILIKE '%' || sqlc.arg(search)::text || '%'
        OR l.country ILIKE '%' || sqlc.arg(search)::text || '%'
        OR l.city ILIKE '%' || sqlc.arg(search)::text || '%'
        OR l.iata ILIKE '%' || sqlc.arg(search)::text || '%'
        OR EXISTS (
            SELECT 1
            FROM location_aliases la
            WHERE la.location_id = l.id
              AND la.alias ILIKE '%' || sqlc.arg(search)::text || '%'
        )
    )
)
SELECT
    id,
    country,
    country_code,
    city,
    name,
    iata,
    is_airport,
    created_at,
    updated_at
FROM matched_locations
ORDER BY
    search_rank ASC,
    CASE WHEN is_airport THEN 0 ELSE 1 END ASC,
    lower(name) ASC
LIMIT 100;


-- name: ListLocationsWithoutAliases :many
SELECT l.id, l.name, l.iata
FROM locations l
WHERE NOT EXISTS (
    SELECT 1
    FROM location_aliases la
    WHERE la.location_id = l.id
)
AND (sqlc.narg('name')::text IS NULL OR l.name ILIKE '%' || sqlc.narg('name')::text || '%')
LIMIT 500;

-- name: ToggleIsAirport :exec
UPDATE locations
SET is_airport = sqlc.arg(is_airport),
    updated_at = now()
WHERE id = sqlc.arg(id);