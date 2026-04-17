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
    driver_title,
    driver_first_name,
    driver_last_name,
    status,
    total_price
FROM reservations
WHERE user_id = @user_id
ORDER BY created_at DESC
LIMIT @query_limit::int
OFFSET @query_offset::int;

-- name: ApplyVoucher :execrows
UPDATE reservations
SET 
    voucher_number = $3,
    vouchered_at = CURRENT_TIMESTAMP
WHERE 
id = $1
AND
user_id = $2;