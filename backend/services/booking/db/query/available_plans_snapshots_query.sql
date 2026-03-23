-- name: InsertAvailablePlansSnapshot :one
INSERT INTO available_plans_snapshots (driver_age, pickup_date, pickup_time, return_date, return_time, country_code, plans)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: DeleteOldAvailablePlansSnapshots :exec
DELETE FROM available_plans_snapshots
where created_at < $1;

-- name: GetSnapshotByID :one
SELECT id, created_at, driver_age, pickup_date, pickup_time, return_date, return_time, country_code, plans
FROM available_plans_snapshots
WHERE id = $1;