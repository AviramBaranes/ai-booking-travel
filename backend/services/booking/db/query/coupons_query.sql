-- name: ListCoupons :many
SELECT id, name, code, discount, is_enabled, created_at, updated_at
FROM coupons
ORDER BY created_at DESC;

-- name: FindCouponByCode :one
SELECT id, name, code, discount, is_enabled, created_at, updated_at
FROM coupons
WHERE code = $1;

-- name: CreateCoupon :one
INSERT INTO coupons (name, code, discount, is_enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, name, code, discount, is_enabled, created_at, updated_at;

-- name: UpdateCoupon :one
UPDATE coupons
SET name       = COALESCE(sqlc.narg(name), name),
    code       = COALESCE(sqlc.narg(code), code),
    discount   = COALESCE(sqlc.narg(discount), discount),
    is_enabled = COALESCE(sqlc.narg(is_enabled), is_enabled),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, name, code, discount, is_enabled, created_at, updated_at;

-- name: DeleteCoupon :exec
DELETE FROM coupons
WHERE id = $1;
