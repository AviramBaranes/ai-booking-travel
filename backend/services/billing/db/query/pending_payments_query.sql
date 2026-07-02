-- name: GetPendingPaymentByID :one
SELECT * FROM pending_customer_payments WHERE id = $1 AND status = 'pending';