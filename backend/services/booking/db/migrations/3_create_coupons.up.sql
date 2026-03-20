CREATE TABLE
    coupons (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        code VARCHAR(100) NOT NULL UNIQUE,
        discount INTEGER NOT NULL CHECK (
            discount > 0
            AND discount <= 100
        ),
        is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_coupons_code ON coupons (code);