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
