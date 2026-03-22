CREATE TABLE
    currencies (
        id SERIAL PRIMARY KEY,
        currency_code VARCHAR(10) NOT NULL UNIQUE,
        currency_iso_name VARCHAR(100) NOT NULL UNIQUE,
        rate NUMERIC(12, 6) NOT NULL CHECK (rate > 0),
        created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
