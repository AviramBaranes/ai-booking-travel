-- name: GetLocationSupplierCode :one
SELECT
  id,
  location_id,
  supplier,
  supplier_location_id,
  enabled,
  created_at,
  updated_at
FROM
  location_supplier_codes
WHERE
  (
    supplier = sqlc.arg (supplier)
    AND supplier_location_id = sqlc.arg (supplier_location_id)
  )
  OR (
    location_id = sqlc.arg (location_id)
    AND supplier = sqlc.arg (supplier)
  )
LIMIT
  1;

-- name: InsertLocationSupplierCode :one
INSERT INTO
  location_supplier_codes (location_id, supplier, supplier_location_id)
VALUES
  (
    sqlc.arg (location_id),
    sqlc.arg (supplier),
    sqlc.arg (supplier_location_id)
  ) RETURNING id,
  location_id,
  supplier,
  supplier_location_id,
  enabled,
  created_at,
  updated_at;