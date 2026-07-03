CREATE TYPE payment_status AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE
    pending_customer_payments (
        id BIGSERIAL PRIMARY KEY,
        user_id BIGINT,
        phone TEXT NOT NULL,
        user_first_name TEXT NOT NULL,
        user_last_name TEXT NOT NULL,
        user_email TEXT NOT NULL,
        snapshot_id BIGINT NOT NULL,
        rate_qualifier TEXT NOT NULL,
        supplier_code TEXT NOT NULL,
        plan_id TEXT NOT NULL,
        include_erp BOOLEAN NOT NULL DEFAULT FALSE,
        selected_addons JSONB,
        driver_title TEXT NOT NULL,
        driver_first_name TEXT NOT NULL,
        driver_last_name TEXT NOT NULL,
        flight_number TEXT,
        status payment_status NOT NULL DEFAULT 'pending',
        reservation_id BIGINT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );