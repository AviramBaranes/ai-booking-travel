-- 
-- 
-- 
-- ---------------- Locations ----------------
-- 
-- 
-- 

CREATE TYPE broker AS ENUM ('flex', 'hertz');

CREATE TABLE
    locations (
        id BIGSERIAL PRIMARY KEY,
        country TEXT NOT NULL,
        country_code TEXT NOT NULL,
        city TEXT,
        name TEXT NOT NULL,
        iata CHAR(3),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        CONSTRAINT locations_iata_uppercase_chk CHECK (
            iata IS NULL
            OR iata = upper(iata)
        )
    );

-- One canonical row per IATA (airports)
CREATE UNIQUE INDEX locations_iata_unique ON locations (iata)
WHERE
    iata IS NOT NULL;

-- Name must be unique within a country
CREATE UNIQUE INDEX locations_country_code_name_unique ON locations (country_code, lower(name));

-- Broker code mappings (how to call each broker for a canonical location)
CREATE TABLE
    location_broker_codes (
        id BIGSERIAL PRIMARY KEY,
        location_id BIGINT NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
        broker broker NOT NULL,
        broker_location_id TEXT NOT NULL,
        enabled BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

-- One row per broker per canonical location
CREATE UNIQUE INDEX uniq_location_broker ON location_broker_codes (location_id, broker);

-- Broker IDs are globally unique per broker
CREATE UNIQUE INDEX uniq_broker_code ON location_broker_codes (broker, broker_location_id);

-- Query helpers
CREATE INDEX idx_location_broker ON location_broker_codes (location_id, broker);

CREATE INDEX idx_broker_lookup ON location_broker_codes (broker, broker_location_id);

-- 
-- 
-- 
-- ---------------- Hertz Prices ----------------
-- 
-- 
-- 
CREATE TABLE
    hertz_markup_rates (
        id BIGSERIAL PRIMARY KEY,
        country TEXT NOT NULL,
        brand TEXT NOT NULL,
        pickup_date_from DATE NOT NULL,
        pickup_date_to DATE NOT NULL,
        car_group TEXT NOT NULL,
        num_of_rental_days_from INT NOT NULL,
        num_of_rental_days_to INT NOT NULL,
        mark_up_gross DOUBLE PRECISION NOT NULL,
        mark_up_net DOUBLE PRECISION NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE INDEX idx_hertz_markup_rates_lookup ON hertz_markup_rates (
    country,
    brand,
    pickup_date_from,
    pickup_date_to,
    car_group,
    num_of_rental_days_from,
    num_of_rental_days_to
);

-- 
-- 
-- 
-- ---------------- Coupons ----------------
-- 
-- 
-- 
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
-- 
-- 
-- 
-- ---------------- Currencies ----------------
-- 
-- 
-- 
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

-- 
-- 
-- 
-- ---------------- Snapshots ----------------
-- 
-- 
-- 
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
-- 
-- 
-- 
-- ---------------- Broker Translations ----------------
-- 
-- 
-- 
CREATE TYPE broker_translation_status AS ENUM ('pending', 'translated', 'verified');

CREATE TABLE
    broker_translations (
        id SERIAL PRIMARY KEY,
        source_text TEXT NOT NULL UNIQUE,
        target_text TEXT,
        status broker_translation_status NOT NULL DEFAULT 'pending',
        confidence_score INTEGER CHECK (
            confidence_score >= 0
            AND confidence_score <= 10
        ),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
-- 
-- 
-- 
-- ---------------- Price Offers ----------------
-- 
-- 
-- 
CREATE TYPE offer_status AS ENUM ('open', 'booked', 'declined');

CREATE TABLE
    price_offers (
        id BIGSERIAL PRIMARY KEY,
        token UUID NOT NULL UNIQUE DEFAULT gen_random_uuid (), 
        agent_id INT NOT NULL, 
        status offer_status NOT NULL DEFAULT 'open',
        name TEXT NOT NULL,

        pickup_location_id TEXT NOT NULL,
        dropoff_location_id TEXT NOT NULL,
        pickup_date DATE NOT NULL,
        return_date DATE NOT NULL,
        pickup_time TEXT NOT NULL,
        dropoff_time TEXT NOT NULL,
        driver_age TEXT NOT NULL,

        supplier_code TEXT NOT NULL,
        car_details JSONB NOT NULL,
        plan_inclusions TEXT[] NOT NULL DEFAULT '{}',
        -- rental_days INT NOT NULL CHECK (rental_days > 0),

        currency_code TEXT NOT NULL,
        purchase_price NUMERIC(12, 2) NOT NULL,
        markup_percentage NUMERIC(12, 2) NOT NULL,
        broker_erp_price NUMERIC(12, 2) NOT NULL,
        bt_erp_price INT NOT NULL,
        total_price INT NOT NULL,
        offered_currency_code TEXT NOT NULL,
        offered_price INT NOT NULL,
        
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE INDEX idx_price_offers_agent_id ON price_offers (agent_id, created_at DESC);

CREATE INDEX idx_price_offers_token ON price_offers (token);