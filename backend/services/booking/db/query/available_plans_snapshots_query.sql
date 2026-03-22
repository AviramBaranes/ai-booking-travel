-- name: InsertAvailablePlansSnapshot :one
INSERT INTO available_plans_snapshots (plans)
VALUES ($1)
RETURNING id;

-- name: DeleteOldAvailablePlansSnapshots :exec
DELETE FROM available_plans_snapshots
where created_at < $1;

-- name: GetSnapshotByID :one
SELECT id, created_at, plans
FROM available_plans_snapshots
WHERE id = $1;