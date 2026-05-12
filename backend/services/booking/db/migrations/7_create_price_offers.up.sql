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