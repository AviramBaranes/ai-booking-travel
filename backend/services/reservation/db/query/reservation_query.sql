-- name: InsertReservation :one
INSERT INTO reservations (
    user_id,
    office_id,
    organization_id,
    is_organization_organic,
    admin_ref_id,
    broker_reservation_id,
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
    dropoff_date,
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
    @office_id,
    @organization_id,
    @is_organization_organic,
    @admin_ref_id,
    @broker_reservation_id,
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
    @dropoff_date,
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
SELECT *
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
    reservation_status,
    total_price
FROM reservations
WHERE user_id = sqlc.arg(user_id)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(name)::VARCHAR IS NULL OR driver_first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR driver_last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_first_name || ' ' || driver_last_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_last_name || ' ' || driver_first_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(pickup_date)::DATE IS NULL OR pickup_date = sqlc.narg(pickup_date)::DATE)
    AND (sqlc.narg(booking_id)::VARCHAR IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(booking_id)::VARCHAR || '%')
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::VARCHAR = 'pickup_date' THEN pickup_date::TIMESTAMP END ASC,
    CASE WHEN sqlc.arg(sort_by)::VARCHAR = 'created_at' OR sqlc.arg(sort_by)::VARCHAR IS NULL THEN created_at END DESC
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountReservationsByUser :one
SELECT COUNT(*)::BIGINT AS total
FROM reservations
WHERE user_id = sqlc.arg(user_id)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(name)::VARCHAR IS NULL OR driver_first_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR driver_last_name ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_first_name || ' ' || driver_last_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%' OR (driver_last_name || ' ' || driver_first_name) ILIKE '%' || sqlc.narg(name)::VARCHAR || '%')
    AND (sqlc.narg(pickup_date)::DATE IS NULL OR pickup_date = sqlc.narg(pickup_date)::DATE)
    AND (sqlc.narg(booking_id)::VARCHAR IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(booking_id)::VARCHAR || '%');

-- name: ApplyVoucher :one
UPDATE reservations
SET 
    reservation_status = 'vouchered',
    voucher_number = $3,
    vouchered_at = CURRENT_TIMESTAMP,
    currency_rate = $4
WHERE 
id = $1
AND
user_id = $2
RETURNING *;

-- name: GetReservationCurrencyCode :one
SELECT currency_code
FROM reservations
WHERE id = $1;

-- name: CancelReservation :exec
UPDATE reservations
SET
    payment_status = CASE
        WHEN payment_status = 'paid' THEN 'refund_pending'
        ELSE payment_status
    END,
    reservation_status = 'canceled',
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1;

-- name: GetPaymentPendingReservations :many
SELECT
    id,
    user_id,
    payment_status,
    driver_title,
    driver_first_name,
    driver_last_name,
    created_at,
    broker_reservation_id,
    vouchered_at,
    voucher_number,
    pickup_date,
    dropoff_date,
    country_code,
    rental_days,
    currency_code,
    purchase_price,
    markup_percentage,
    bt_erp_price,
    broker_erp_price,
    total_price
FROM reservations
WHERE
    (reservation_status = 'vouchered' AND payment_status = 'unpaid')
OR
    (reservation_status = 'canceled' AND payment_status = 'refund_pending');

-- name: GetPaymentPendingReservationsByBillingEntity :many
SELECT
    id,
    broker_reservation_id,
    payment_status,
    reservation_status,
    purchase_price,
    markup_percentage,
    bt_erp_price,
    broker_erp_price,
    total_price,
    currency_code,
    created_at,
    pickup_date
FROM reservations
WHERE
    (
        (
            sqlc.narg(office_id)::BIGINT IS NULL 
            AND  (sqlc.narg(organization_id)::BIGINT) IS NOT NULL 
            AND organization_id = sqlc.narg(organization_id)::BIGINT
            AND is_organization_organic = TRUE
        )
    OR
        (
            sqlc.narg(organization_id)::BIGINT IS NULL 
            AND  (sqlc.narg(office_id)::BIGINT) IS NOT NULL 
            AND office_id = sqlc.narg(office_id)::BIGINT
            AND is_organization_organic = FALSE
        )
    )
AND(
    (reservation_status = 'vouchered' AND payment_status = 'unpaid')
OR
    (reservation_status = 'canceled' AND payment_status = 'refund_pending'));

-- name: ListBusinessesBalancesReport :many
WITH billing_reservations AS (
    SELECT
        CASE
            WHEN is_organization_organic = TRUE THEN 'organization'
            ELSE 'office'
        END::TEXT AS billing_entity_type,
        CASE
            WHEN is_organization_organic = TRUE THEN organization_id
            ELSE office_id
        END::BIGINT AS billing_entity_id,
        reservation_status,
        payment_status,
        (total_price::NUMERIC * currency_rate)::DOUBLE PRECISION AS balance
    FROM reservations
    WHERE
        (
            is_organization_organic = TRUE
            AND organization_id IS NOT NULL
        )
    OR
        (
            is_organization_organic = FALSE
            AND office_id IS NOT NULL
        )
)
SELECT
    billing_entity_type,
    billing_entity_id,
    COUNT(*) FILTER (WHERE reservation_status = 'booked')::BIGINT AS open_reservations_count,
    COALESCE(SUM(balance) FILTER (WHERE reservation_status = 'booked'), 0)::DOUBLE PRECISION AS total_open_balance,
    COUNT(*) FILTER (WHERE reservation_status = 'vouchered' AND payment_status = 'unpaid')::BIGINT AS payment_pending_reservations_count,
    COALESCE(SUM(balance) FILTER (WHERE reservation_status = 'vouchered' AND payment_status = 'unpaid'), 0)::DOUBLE PRECISION AS total_payment_pending_balance,
    COUNT(*) FILTER (WHERE reservation_status = 'canceled' AND payment_status = 'refund_pending')::BIGINT AS refund_pending_reservations_count,
    COALESCE(SUM(balance) FILTER (WHERE reservation_status = 'canceled' AND payment_status = 'refund_pending'), 0)::DOUBLE PRECISION AS total_refund_pending_balance,
    (
        COALESCE(SUM(balance) FILTER (WHERE reservation_status = 'booked' OR (reservation_status = 'vouchered' AND payment_status = 'unpaid')), 0)::DOUBLE PRECISION
            - COALESCE(SUM(balance) FILTER (WHERE reservation_status = 'canceled' AND payment_status = 'refund_pending'), 0)::DOUBLE PRECISION
    )::DOUBLE PRECISION AS total_balance
FROM billing_reservations
GROUP BY billing_entity_type, billing_entity_id
HAVING
    COUNT(*) FILTER (WHERE reservation_status = 'booked') > 0
OR
    COUNT(*) FILTER (WHERE reservation_status = 'vouchered' AND payment_status = 'unpaid') > 0
OR
    COUNT(*) FILTER (WHERE reservation_status = 'canceled' AND payment_status = 'refund_pending') > 0
ORDER BY billing_entity_type, billing_entity_id;

-- name: ResolveReservationsPayment :exec
UPDATE reservations
SET payment_status = CASE payment_status
    WHEN 'unpaid' THEN 'paid'
    WHEN 'refund_pending' THEN 'refunded'
    ELSE payment_status
END,
updated_at = CURRENT_TIMESTAMP
WHERE id = ANY(@ids::BIGINT[])
AND payment_status IN ('unpaid', 'refund_pending');

-- name: ListReservationsReport :many
SELECT *
FROM reservations
WHERE
    (sqlc.narg(pickup_date_from)::DATE IS NULL OR pickup_date >= sqlc.narg(pickup_date_from)::DATE)
    AND (sqlc.narg(pickup_date_to)::DATE IS NULL OR pickup_date <= sqlc.narg(pickup_date_to)::DATE)
    AND (sqlc.narg(created_date_from)::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg(created_date_from)::TIMESTAMPTZ)
    AND (sqlc.narg(created_date_to)::TIMESTAMPTZ IS NULL OR created_at <= sqlc.narg(created_date_to)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_from)::TIMESTAMPTZ IS NULL OR vouchered_at >= sqlc.narg(vouchered_at_from)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_to)::TIMESTAMPTZ IS NULL OR vouchered_at <= sqlc.narg(vouchered_at_to)::TIMESTAMPTZ)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(broker)::broker IS NULL OR broker = sqlc.narg(broker)::broker)
    AND (sqlc.narg(supplier_code)::TEXT IS NULL OR supplier_code = sqlc.narg(supplier_code)::TEXT)
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR organization_id = sqlc.narg(organization_id)::BIGINT)
    AND (sqlc.narg(office_id)::BIGINT IS NULL OR office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(agent_id)::BIGINT IS NULL OR user_id = sqlc.narg(agent_id)::BIGINT)
    AND (NOT sqlc.arg(is_business)::BOOLEAN OR (office_id IS NOT NULL AND organization_id IS NOT NULL))
ORDER BY created_at DESC
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountReservationsReport :one
SELECT COUNT(*)::BIGINT AS total
FROM reservations
WHERE
    (sqlc.narg(pickup_date_from)::DATE IS NULL OR pickup_date >= sqlc.narg(pickup_date_from)::DATE)
    AND (sqlc.narg(pickup_date_to)::DATE IS NULL OR pickup_date <= sqlc.narg(pickup_date_to)::DATE)
    AND (sqlc.narg(created_date_from)::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg(created_date_from)::TIMESTAMPTZ)
    AND (sqlc.narg(created_date_to)::TIMESTAMPTZ IS NULL OR created_at <= sqlc.narg(created_date_to)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_from)::TIMESTAMPTZ IS NULL OR vouchered_at >= sqlc.narg(vouchered_at_from)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_to)::TIMESTAMPTZ IS NULL OR vouchered_at <= sqlc.narg(vouchered_at_to)::TIMESTAMPTZ)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(broker)::broker IS NULL OR broker = sqlc.narg(broker)::broker)
    AND (sqlc.narg(supplier_code)::TEXT IS NULL OR supplier_code = sqlc.narg(supplier_code)::TEXT)
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR organization_id = sqlc.narg(organization_id)::BIGINT)
    AND (sqlc.narg(office_id)::BIGINT IS NULL OR office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(agent_id)::BIGINT IS NULL OR user_id = sqlc.narg(agent_id)::BIGINT)
    AND (NOT sqlc.arg(is_business)::BOOLEAN OR (office_id IS NOT NULL AND organization_id IS NOT NULL));

-- name: GetOpenReservationsPickingUpWithinWeek :many
SELECT
    id,
    user_id,
    broker_reservation_id,
    driver_title,
    driver_first_name,
    driver_last_name,
    pickup_date,
    pickup_time
FROM reservations
WHERE
    reservation_status = 'booked'
    AND pickup_date <= CURRENT_DATE + INTERVAL '7 days';