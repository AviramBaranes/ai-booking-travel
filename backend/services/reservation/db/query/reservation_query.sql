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
    price_before_discount,
    price_after_discount,
    discount_percentage,
    erp_price,
    total_price,
    pickup_date,
    return_date,
    rental_days,
    driver_title,
    driver_first_name,
    driver_last_name,
    driver_age
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
    @price_before_discount,
    @price_after_discount,
    @discount_percentage,
    @erp_price,
    @total_price,
    @pickup_date,
    @return_date,
    @rental_days,
    @driver_title,
    @driver_first_name,
    @driver_last_name,
    @driver_age
) RETURNING id;

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