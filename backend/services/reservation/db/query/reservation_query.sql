-- name: InsertReservation :one
INSERT INTO reservations (
    user_id,
    broker_reservation_id,
    status,
    broker,
    supplier_code,
    car_details,
    plan_inclusions,
    country_code,
    currency_code,
    currency_rate,
    purchase_price,
    markup_percentage,
    discount_percentage,
    broker_erp_price,
    bt_erp_price,
    vat_percentage,
    total_price,
    pickup_date,
    return_date,
    pickup_time,
    dropoff_time,
    rental_days,
    driver_title,
    driver_first_name,
    driver_last_name,
    driver_age,
    pickup_location_name,
    dropoff_location_name
) VALUES (
    @user_id,
    @broker_reservation_id,
    @status,
    @broker,
    @supplier_code,
    @car_details,
    @plan_inclusions,
    @country_code,
    @currency_code,
    @currency_rate,
    @purchase_price,
    @markup_percentage,
    @discount_percentage,
    @broker_erp_price,
    @bt_erp_price,
    @vat_percentage,
    @total_price,
    @pickup_date,
    @return_date,
    @pickup_time,
    @dropoff_time,
    @rental_days,
    @driver_title,
    @driver_first_name,
    @driver_last_name,
    @driver_age,
    @pickup_location_name,
    @dropoff_location_name
) RETURNING id;

-- name: GetReservationByID :one
SELECT
    id,
    user_id,
    broker_reservation_id,
    status,
    broker,
    supplier_code,
    car_details,
    plan_inclusions,
    country_code,
    currency_code,
    currency_rate,
    purchase_price,
    markup_percentage,
    discount_percentage,
    broker_erp_price,
    bt_erp_price,
    vat_percentage,
    total_price,
    pickup_date,
    return_date,
    pickup_time,
    dropoff_time,
    rental_days,
    driver_title,
    driver_first_name,
    driver_last_name,
    driver_age,
    pickup_location_name,
    dropoff_location_name,
    voucher_number,
    vouchered_at,
    created_at
FROM reservations
WHERE id = @id;

-- name: ListReservationsByUser :many
SELECT
    id,
    broker_reservation_id,
    created_at,
    country_code,
    pickup_date,
    pickup_location_name,
    driver_title,
    driver_first_name,
    driver_last_name,
    status,
    total_price
FROM reservations
WHERE user_id = sqlc.arg(user_id)
    AND (sqlc.narg(status)::reservation_status IS NULL OR status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(name)::VARCHAR IS NULL OR driver_first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR driver_last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_first_name || ' ' || driver_last_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_last_name || ' ' || driver_first_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(pickup_date)::DATE IS NULL OR pickup_date = sqlc.narg(pickup_date)::DATE)
    AND (sqlc.narg(booking_id)::VARCHAR IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(booking_id)::VARCHAR || '%')
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::VARCHAR = 'pickup_date' THEN pickup_date::TIMESTAMP END DESC,
    CASE WHEN sqlc.arg(sort_by)::VARCHAR = 'created_at' OR sqlc.arg(sort_by)::VARCHAR IS NULL THEN created_at END DESC
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountReservationsByUser :one
SELECT COUNT(*)::BIGINT AS total
FROM reservations
WHERE user_id = sqlc.arg(user_id)
    AND (sqlc.narg(status)::reservation_status IS NULL OR status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(name)::VARCHAR IS NULL OR driver_first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR driver_last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_first_name || ' ' || driver_last_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_last_name || ' ' || driver_first_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(pickup_date)::DATE IS NULL OR pickup_date = sqlc.narg(pickup_date)::DATE)
    AND (sqlc.narg(booking_id)::VARCHAR IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(booking_id)::VARCHAR || '%');

-- name: ApplyVoucher :execrows
UPDATE reservations
SET 
    status = 'vouchered',
    voucher_number = $3,
    vouchered_at = CURRENT_TIMESTAMP
WHERE 
id = $1
AND
user_id = $2;

-- name: CancelReservation :exec
UPDATE reservations
SET
    status = 'canceled',
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1;