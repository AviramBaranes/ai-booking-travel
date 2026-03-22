-- name: ListCurrencies :many
SELECT id, currency_code, currency_iso_name, rate, created_at, updated_at
FROM currencies
ORDER BY created_at DESC;

-- name: FindCurrencyByISOName :one
SELECT id, currency_code, currency_iso_name, rate, created_at, updated_at
FROM currencies
WHERE currency_iso_name = $1;

-- name: CreateCurrency :one
INSERT INTO currencies (currency_code, currency_iso_name, rate, created_at, updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, currency_code, currency_iso_name, rate, created_at, updated_at;

-- name: UpdateCurrency :one
UPDATE currencies
SET currency_code     = COALESCE(sqlc.narg(currency_code), currency_code),
    currency_iso_name = COALESCE(sqlc.narg(currency_iso_name), currency_iso_name),
    rate              = COALESCE(sqlc.narg(rate), rate),
    updated_at        = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, currency_code, currency_iso_name, rate, created_at, updated_at;

-- name: DeleteCurrency :exec
DELETE FROM currencies
WHERE id = $1;
