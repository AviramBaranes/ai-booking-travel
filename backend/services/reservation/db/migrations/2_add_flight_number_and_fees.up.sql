ALTER TABLE reservations
ADD COLUMN flight_number TEXT,
ADD COLUMN pay_at_pickup JSONB NOT NULL DEFAULT '{}'::jsonb;