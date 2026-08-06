-- name: GetSupplierTermsByID :one
SELECT terms
FROM supplier_terms
WHERE id = $1;

-- name: UpsertSupplierTerms :one
INSERT INTO supplier_terms (
    broker,
    supplier_code,
    pickup_location_id,
    terms
) VALUES (
    @broker,
    @supplier_code,
    @pickup_location_id,
    @terms
)
ON CONFLICT (broker, supplier_code, pickup_location_id)
DO UPDATE SET terms = EXCLUDED.terms
RETURNING id;
