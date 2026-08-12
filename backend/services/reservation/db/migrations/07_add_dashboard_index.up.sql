-- The admin dashboard scans a created_at window across all users. Every existing index
-- is a (user_id, ...) composite, so that scan has no index to use.
CREATE INDEX IF NOT EXISTS idx_reservations_created_at ON reservations (created_at);
