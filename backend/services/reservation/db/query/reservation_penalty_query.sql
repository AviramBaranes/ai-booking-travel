-- name: InsertReservationPenalty :one
INSERT INTO reservation_penalties (
    reservation_id,
    penalty_type,
    created_by_user_id,
    currency_code,
    currency_rate,
    amount
) VALUES (
    @reservation_id,
    @penalty_type,
    @created_by_user_id,
    @currency_code,
    @currency_rate,
    @amount
) RETURNING *;

-- name: GetReservationPenaltyByReservationID :one
SELECT *
FROM reservation_penalties
WHERE reservation_id = @reservation_id;

-- name: GetPaymentPendingPenalties :many
SELECT
    p.id,
    p.reservation_id,
    p.penalty_type,
    p.amount,
    p.currency_code,
    r.broker_reservation_id,
    r.user_id,
    r.created_at,
    r.vouchered_at,
    r.voucher_number,
    r.driver_title,
    r.driver_first_name,
    r.driver_last_name,
    r.pickup_date,
    r.dropoff_date,
    r.rental_days,
    r.country_code
FROM reservation_penalties p
JOIN reservations r ON r.id = p.reservation_id
WHERE p.paid_at IS NULL;

-- name: ResolvePenaltiesPayment :exec
UPDATE reservation_penalties
SET
    paid_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ANY(@ids::BIGINT[])
AND paid_at IS NULL;

-- name: GetPaymentPendingPenaltiesByBillingEntity :many
SELECT
    p.id,
    p.reservation_id,
    r.broker_reservation_id,
    p.penalty_type,
    p.amount,
    p.currency_code,
    p.currency_rate,
    p.created_at
FROM reservation_penalties p
JOIN reservations r ON r.id = p.reservation_id
WHERE
    (
        (
            sqlc.narg(office_id)::BIGINT IS NULL
            AND (sqlc.narg(organization_id)::BIGINT) IS NOT NULL
            AND r.organization_id = sqlc.narg(organization_id)::BIGINT
            AND r.is_organization_organic = TRUE
        )
    OR
        (
            sqlc.narg(organization_id)::BIGINT IS NULL
            AND (sqlc.narg(office_id)::BIGINT) IS NOT NULL
            AND r.office_id = sqlc.narg(office_id)::BIGINT
            AND r.is_organization_organic = FALSE
        )
    )
    AND p.paid_at IS NULL
ORDER BY p.created_at;
