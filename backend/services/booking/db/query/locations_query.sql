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

-- name: SearchLocations :many
SELECT *
FROM locations
WHERE EXISTS (
    SELECT 1 
    FROM location_broker_codes lbc 
    WHERE lbc.location_id = locations.id 
      AND lbc.enabled = TRUE
  )
  AND (
    name ILIKE '%' || sqlc.arg(search)::text || '%'
    OR country ILIKE '%' || sqlc.arg(search)::text || '%'
    OR iata ILIKE '%' || sqlc.arg(search)::text || '%'
    OR city ILIKE '%' || sqlc.arg(search)::text || '%'
  )
ORDER BY iata ASC
LIMIT 30;

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

-- name: DeleteLocationByID :exec
DELETE FROM locations
WHERE id = sqlc.arg(id);