CREATE UNLOGGED TABLE
    available_plans_snapshots (
        id BIGSERIAL PRIMARY KEY,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        driver_age TEXT NOT NULL,
        pickup_date DATE NOT NULL,
        pickup_time TEXT NOT NULL,
        return_date DATE NOT NULL,
        return_time TEXT NOT NULL,
        country_code TEXT NOT NULL,
        plans JSON NOT NULL
    );

CREATE INDEX available_plans_snapshots_created_at_idx ON available_plans_snapshots (created_at);