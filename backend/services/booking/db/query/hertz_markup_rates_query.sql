-- name: GetHertzMarkupRates :many
-- Used by the markup logic to fetch rates for a given search.
SELECT car_group, brand, mark_up_gross, mark_up_net
FROM hertz_markup_rates
WHERE country = $1
  AND pickup_date_from <= sqlc.arg(pickup_date)::date
  AND pickup_date_to >= sqlc.arg(pickup_date)::date
  AND num_of_rental_days_from <= sqlc.arg(rental_days)::int
  AND num_of_rental_days_to >= sqlc.arg(rental_days)::int
  AND car_group = ANY(sqlc.arg(car_groups)::text[]);

-- name: CountHertzMarkupRates :one
-- Count total rows matching the same filters (for pagination).
SELECT COUNT(*) AS total
FROM hertz_markup_rates
WHERE (sqlc.narg(country)::text IS NULL OR country ILIKE '%' || sqlc.narg(country) || '%')
  AND (sqlc.narg(brand)::text IS NULL OR brand ILIKE '%' || sqlc.narg(brand) || '%')
  AND (sqlc.narg(car_group)::text IS NULL OR car_group ILIKE '%' || sqlc.narg(car_group) || '%');

-- name: ListHertzMarkupRates :many
-- Admin listing with pagination, optional filtering, and sorting.
SELECT *
FROM hertz_markup_rates
WHERE (sqlc.narg(country)::text IS NULL OR country ILIKE '%' || sqlc.narg(country) || '%')
  AND (sqlc.narg(brand)::text IS NULL OR brand ILIKE '%' || sqlc.narg(brand) || '%')
  AND (sqlc.narg(car_group)::text IS NULL OR car_group ILIKE '%' || sqlc.narg(car_group) || '%')
ORDER BY
  CASE WHEN sqlc.arg(sort_field)::text = 'country' AND sqlc.arg(sort_dir)::text = 'asc' THEN country END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'country' AND sqlc.arg(sort_dir)::text = 'desc' THEN country END DESC,
  CASE WHEN sqlc.arg(sort_field)::text = 'brand' AND sqlc.arg(sort_dir)::text = 'asc' THEN brand END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'brand' AND sqlc.arg(sort_dir)::text = 'desc' THEN brand END DESC,
  CASE WHEN sqlc.arg(sort_field)::text = 'car_group' AND sqlc.arg(sort_dir)::text = 'asc' THEN car_group END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'car_group' AND sqlc.arg(sort_dir)::text = 'desc' THEN car_group END DESC,
  CASE WHEN sqlc.arg(sort_field)::text = 'pickup_date_from' AND sqlc.arg(sort_dir)::text = 'asc' THEN pickup_date_from END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'pickup_date_from' AND sqlc.arg(sort_dir)::text = 'desc' THEN pickup_date_from END DESC,
  CASE WHEN sqlc.arg(sort_field)::text = 'num_of_rental_days_from' AND sqlc.arg(sort_dir)::text = 'asc' THEN num_of_rental_days_from END ASC,
  CASE WHEN sqlc.arg(sort_field)::text = 'num_of_rental_days_from' AND sqlc.arg(sort_dir)::text = 'desc' THEN num_of_rental_days_from END DESC,
  id ASC
LIMIT sqlc.arg(query_limit)::int
OFFSET sqlc.arg(query_offset)::int;

-- name: InsertHertzMarkupRate :one
INSERT INTO hertz_markup_rates (
    country, brand, pickup_date_from, pickup_date_to,
    car_group, num_of_rental_days_from, num_of_rental_days_to,
    mark_up_gross, mark_up_net
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateHertzMarkupRate :one
UPDATE hertz_markup_rates
SET country = $2,
    brand = $3,
    pickup_date_from = $4,
    pickup_date_to = $5,
    car_group = $6,
    num_of_rental_days_from = $7,
    num_of_rental_days_to = $8,
    mark_up_gross = $9,
    mark_up_net = $10,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteHertzMarkupRate :one
DELETE FROM hertz_markup_rates
WHERE id = $1
RETURNING id;
