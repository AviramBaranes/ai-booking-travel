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
