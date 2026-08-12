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
    coupon_name,
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
    dropoff_location_name,
    flight_number,
    pay_at_pickup,
    excess,
    excess_currency,
    pickup_location_code,
    dropoff_location_code,
    supplier_terms_id,
    pickup_details,
    dropoff_details
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
    @coupon_name,
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
    @dropoff_location_name,
    @flight_number,
    @pay_at_pickup,
    @excess,
    @excess_currency,
    @pickup_location_code,
    @dropoff_location_code,
    @supplier_terms_id,
    @pickup_details,
    @dropoff_details
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

-- name: ApplyVoucher :exec
UPDATE reservations
SET 
    reservation_status = 'vouchered',
    voucher_number = $3,
    vouchered_at = CURRENT_TIMESTAMP,
    currency_rate = $4
WHERE 
id = $1
AND
user_id = $2;

-- name: CancelReservation :exec
UPDATE reservations
SET
    payment_status = CASE
        WHEN payment_status = 'paid' THEN 
            CASE WHEN payment_doc_num IS NOT NULL THEN 'refunded'::payment_status
            ELSE 'refund_pending'::payment_status
            END
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
    discount_percentage,
    bt_erp_price,
    broker_erp_price,
    total_price,
    currency_code,
    currency_rate,
    created_at,
    vouchered_at,
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
    (reservation_status = 'canceled' AND payment_status = 'refund_pending'))
ORDER BY vouchered_at;

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
        CASE
            WHEN reservation_status = 'canceled' AND payment_status = 'refund_pending' THEN -total_price
            ELSE total_price
        END AS total_price,
        currency_code,
        currency_rate
    FROM reservations
    WHERE
        (
            (is_organization_organic = TRUE AND organization_id IS NOT NULL)
            OR
            (is_organization_organic = FALSE AND office_id IS NOT NULL)
        )
        AND (
            (reservation_status = 'vouchered' AND payment_status = 'unpaid')
            OR
            (reservation_status = 'canceled' AND payment_status = 'refund_pending')
        )
)
SELECT
    billing_entity_type,
    billing_entity_id,
    COALESCE(SUM(total_price) FILTER (WHERE currency_code = 'EUR'), 0)::DOUBLE PRECISION AS total_eur,
    COALESCE(SUM(total_price) FILTER (WHERE currency_code = 'USD'), 0)::DOUBLE PRECISION AS total_usd,
    COALESCE(SUM(total_price * currency_rate) FILTER (WHERE currency_code NOT IN ('EUR', 'USD')), 0)::DOUBLE PRECISION AS total_other_converted
FROM billing_reservations
GROUP BY billing_entity_type, billing_entity_id
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
    (sqlc.narg(broker_reservation_id)::TEXT IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(broker_reservation_id)::TEXT || '%')
    AND(sqlc.narg(pickup_date_from)::DATE IS NULL OR pickup_date >= sqlc.narg(pickup_date_from)::DATE)
    AND (sqlc.narg(pickup_date_to)::DATE IS NULL OR pickup_date <= sqlc.narg(pickup_date_to)::DATE)
    AND (sqlc.narg(created_date_from)::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg(created_date_from)::TIMESTAMPTZ)
    AND (sqlc.narg(created_date_to)::TIMESTAMPTZ IS NULL OR created_at < sqlc.narg(created_date_to)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_from)::TIMESTAMPTZ IS NULL OR vouchered_at >= sqlc.narg(vouchered_at_from)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_to)::TIMESTAMPTZ IS NULL OR vouchered_at < sqlc.narg(vouchered_at_to)::TIMESTAMPTZ)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(broker)::broker IS NULL OR broker = sqlc.narg(broker)::broker)
    AND (COALESCE(cardinality(sqlc.arg(supplier_codes)::TEXT[]), 0) = 0 OR supplier_code = ANY(sqlc.arg(supplier_codes)::TEXT[]))
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR organization_id = sqlc.narg(organization_id)::BIGINT)
    AND (sqlc.narg(office_id)::BIGINT IS NULL OR office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(agent_id)::BIGINT IS NULL OR user_id = sqlc.narg(agent_id)::BIGINT)
    AND (sqlc.narg(is_business)::BOOLEAN IS NULL OR (sqlc.narg(is_business)::BOOLEAN = TRUE AND office_id IS NOT NULL AND organization_id IS NOT NULL) OR (sqlc.narg(is_business)::BOOLEAN = FALSE AND office_id IS NULL AND organization_id IS NULL))
    AND (NOT sqlc.arg(skip_canceled)::BOOLEAN OR reservation_status != 'canceled')
ORDER BY created_at DESC
LIMIT  sqlc.arg(page_size)::BIGINT
OFFSET sqlc.arg(page_offset)::BIGINT;

-- name: CountReservationsReport :one
SELECT 
COUNT(*)::BIGINT AS count,
COALESCE(SUM(total_price * currency_rate), 0)::DOUBLE PRECISION AS total_sales,
COALESCE(SUM(purchase_price * currency_rate), 0)::DOUBLE PRECISION AS total_car_cost,
COALESCE(SUM(broker_erp_price * currency_rate), 0)::DOUBLE PRECISION AS total_broker_erp_cost
FROM reservations
WHERE
    (sqlc.narg(broker_reservation_id)::TEXT IS NULL OR broker_reservation_id ILIKE '%' || sqlc.narg(broker_reservation_id)::TEXT || '%')
    AND (sqlc.narg(pickup_date_from)::DATE IS NULL OR pickup_date >= sqlc.narg(pickup_date_from)::DATE)
    AND (sqlc.narg(pickup_date_to)::DATE IS NULL OR pickup_date <= sqlc.narg(pickup_date_to)::DATE)
    AND (sqlc.narg(created_date_from)::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg(created_date_from)::TIMESTAMPTZ)
    AND (sqlc.narg(created_date_to)::TIMESTAMPTZ IS NULL OR created_at < sqlc.narg(created_date_to)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_from)::TIMESTAMPTZ IS NULL OR vouchered_at >= sqlc.narg(vouchered_at_from)::TIMESTAMPTZ)
    AND (sqlc.narg(vouchered_at_to)::TIMESTAMPTZ IS NULL OR vouchered_at < sqlc.narg(vouchered_at_to)::TIMESTAMPTZ)
    AND (sqlc.narg(status)::reservation_status IS NULL OR reservation_status = sqlc.narg(status)::reservation_status)
    AND (sqlc.narg(broker)::broker IS NULL OR broker = sqlc.narg(broker)::broker)
    AND (COALESCE(cardinality(sqlc.arg(supplier_codes)::TEXT[]), 0) = 0 OR supplier_code = ANY(sqlc.arg(supplier_codes)::TEXT[]))
    AND (sqlc.narg(organization_id)::BIGINT IS NULL OR organization_id = sqlc.narg(organization_id)::BIGINT)
    AND (sqlc.narg(office_id)::BIGINT IS NULL OR office_id = sqlc.narg(office_id)::BIGINT)
    AND (sqlc.narg(agent_id)::BIGINT IS NULL OR user_id = sqlc.narg(agent_id)::BIGINT)
    AND (sqlc.narg(is_business)::BOOLEAN IS NULL OR (sqlc.narg(is_business)::BOOLEAN = TRUE AND office_id IS NOT NULL AND organization_id IS NOT NULL) OR (sqlc.narg(is_business)::BOOLEAN = FALSE AND office_id IS NULL AND organization_id IS NULL))
    AND (NOT sqlc.arg(skip_canceled)::BOOLEAN OR reservation_status != 'canceled');

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

-- name: VoucherReservationAfterPayment :one
UPDATE reservations
SET
    reservation_status = 'vouchered',
    payment_status = 'paid',
    payment_confirmation_code = $2,
    payment_doc_num = $3,
    vouchered_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
AND
    reservation_status = 'booked'
AND
    payment_status = 'unpaid'
RETURNING *;

-- name: UpdateReservationCurrencyRate :exec
UPDATE reservations
SET
    currency_rate = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1;

-- name: GetBillingReservation :one
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
    id = $1;

-- name: SaveInvoiceDocNum :exec
UPDATE reservations
SET invoice_doc_num = $2
WHERE id = $1;

-- name: ListUnpaidSupplierReservations :many
SELECT
    id,
    broker_reservation_id,
    driver_title,
    driver_first_name,
    driver_last_name,
    pickup_date,
    pickup_location_name,
    rental_days,
    total_price,
    purchase_price,
    broker_erp_price,
    currency_code,
    reservation_status,
    payment_status
FROM reservations
WHERE
    broker = sqlc.arg(broker)::broker
    AND reservation_status != 'canceled'
    AND supplier_paid_at IS NULL
ORDER BY pickup_date, id;

-- name: ListUnpaidSupplierReservationsByIDs :many
SELECT
    id,
    broker_reservation_id,
    broker,
    purchase_price,
    broker_erp_price,
    currency_code,
    currency_rate
FROM reservations
WHERE
    id = ANY(sqlc.arg(ids)::BIGINT[])
    AND reservation_status != 'canceled'
    AND supplier_paid_at IS NULL;

-- name: ListReservationsByBrokerReservationIDs :many
SELECT
    id,
    broker_reservation_id,
    reservation_status,
    supplier_paid_at
FROM reservations
WHERE
    broker = sqlc.arg(broker)::broker
    AND broker_reservation_id = ANY(sqlc.arg(broker_reservation_ids)::TEXT[]);

-- name: MarkReservationsPaidToSupplier :exec
UPDATE reservations
SET
    supplier_paid_at = sqlc.arg(supplier_paid_at)::TIMESTAMPTZ,
    supplier_expense_id = sqlc.narg(supplier_expense_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ANY(sqlc.arg(ids)::BIGINT[])
AND supplier_paid_at IS NULL;