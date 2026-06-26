-- name: GetBrokerMarkupRateByCountryCode :one
-- Used by the markup logic to fetch rates for a given search.
SELECT mark_up_gross, mark_up_net
FROM markup_rates
WHERE country_code = $1
  AND broker = sqlc.arg(broker)::broker;

-- name: CountMarkupRates :one
-- Count total rows matching the same filters (for pagination).
SELECT COUNT(*) AS total
FROM markup_rates
WHERE (sqlc.narg(country_code)::text IS NULL OR country_code ILIKE '%' || sqlc.narg(country_code) || '%')
  AND (sqlc.narg(broker)::text IS NULL OR broker::text ILIKE '%' || sqlc.narg(broker) || '%');

-- name: ListMarkupRates :many
-- Admin listing with pagination, optional filtering, and sorting.
SELECT *
FROM markup_rates
WHERE (sqlc.narg(country_code)::text IS NULL OR country_code ILIKE '%' || sqlc.narg(country_code) || '%')
  AND (sqlc.narg(broker)::text IS NULL OR broker::text ILIKE '%' || sqlc.narg(broker) || '%')
ORDER BY
  CASE WHEN sqlc.arg(sort_field)::text = 'country_code' AND sqlc.arg(sort_dir)::text = 'asc' THEN country_code END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'country_code' AND sqlc.arg(sort_dir)::text = 'desc' THEN country_code END DESC,
  CASE WHEN sqlc.arg(sort_field)::text = 'broker' AND sqlc.arg(sort_dir)::text = 'asc' THEN broker::text END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'broker' AND sqlc.arg(sort_dir)::text = 'desc' THEN broker::text END DESC,
  id ASC
LIMIT sqlc.arg(query_limit)::int
OFFSET sqlc.arg(query_offset)::int;

-- name: InsertMarkupRate :one
INSERT INTO markup_rates (country_code, broker, mark_up_gross, mark_up_net)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMarkupRate :one
UPDATE markup_rates
SET country_code = $2,
    broker = $3,
    mark_up_gross = $4,
    mark_up_net = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteMarkupRate :one
DELETE FROM markup_rates
WHERE id = $1
RETURNING id;
