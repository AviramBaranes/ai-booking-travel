-- name: InsertAvailablePlansSnapshot :one
INSERT INTO available_plans_snapshots (driver_age, pickup_date, pickup_time, dropoff_date, dropoff_time, country_code, plans, suppliers_info)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: DeleteOldAvailablePlansSnapshots :exec
DELETE FROM available_plans_snapshots
where created_at < $1;

-- name: GetSnapshotByID :one
SELECT id, created_at, driver_age, pickup_date, pickup_time, dropoff_date, dropoff_time, country_code, plans, suppliers_info
FROM available_plans_snapshots
WHERE id = $1;

-- name: DeleteSnapshotByID :exec
DELETE FROM available_plans_snapshots WHERE id = $1;