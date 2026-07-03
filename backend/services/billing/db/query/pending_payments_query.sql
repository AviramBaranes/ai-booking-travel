-- name: GetPendingPaymentByID :one
SELECT * FROM pending_customer_payments WHERE id = $1 AND status = 'pending';

-- name: CreatePendingPayment :one
INSERT INTO pending_customer_payments (
    user_id,
    phone,
    user_first_name,
    user_last_name,
    user_email,
    snapshot_id,
    rate_qualifier,
    supplier_code,
    plan_id,
    include_erp,
    selected_addons,
    driver_title,
    driver_first_name,
    driver_last_name,
    flight_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING id;