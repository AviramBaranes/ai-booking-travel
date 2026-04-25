CREATE TYPE reservation_status AS ENUM ('booked', 'vouchered', 'canceled');

CREATE TYPE payment_status AS ENUM ('unpaid', 'paid', 'refund_pending', 'refunded');

CREATE TYPE broker AS ENUM ('flex', 'hertz');

CREATE TABLE
    reservations (
        id BIGSERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        broker_reservation_id TEXT NOT NULL,
        reservation_status reservation_status NOT NULL DEFAULT 'booked',
        payment_status payment_status NOT NULL DEFAULT 'unpaid',
        broker broker NOT NULL,
        supplier_code TEXT NOT NULL,
        car_details JSONB NOT NULL,
        plan_inclusions TEXT[] NOT NULL DEFAULT '{}',
        country_code TEXT NOT NULL,
        currency_code TEXT NOT NULL,
        currency_rate NUMERIC(12, 4) NOT NULL,
        purchase_price NUMERIC(12, 2) NOT NULL,
        markup_percentage NUMERIC(12, 2) NOT NULL,
        broker_erp_price NUMERIC(12, 2) NOT NULL,
        discount_percentage INT NOT NULL DEFAULT 0,
        bt_erp_price INT NOT NULL DEFAULT 0,
        vat_percentage NUMERIC(5, 2) NOT NULL,
        total_price INT NOT NULL,
        pickup_date DATE NOT NULL,
        return_date DATE NOT NULL,
        pickup_time TEXT NOT NULL DEFAULT '',
        dropoff_time TEXT NOT NULL DEFAULT '',
        rental_days INT NOT NULL CHECK (rental_days > 0),
        driver_title TEXT NOT NULL,
        driver_first_name TEXT NOT NULL,
        driver_last_name TEXT NOT NULL,
        driver_age INT NOT NULL CHECK (driver_age >= 18),
        pickup_location_name TEXT NOT NULL DEFAULT '',
        dropoff_location_name TEXT NOT NULL DEFAULT '',
        voucher_number TEXT,
        vouchered_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );

CREATE INDEX idx_reservations_user_id_created_at ON reservations (user_id, created_at DESC);

CREATE INDEX idx_reservations_user_id_reservation_status_created_at ON reservations (user_id, reservation_status, created_at DESC);

CREATE INDEX idx_reservations_user_id_pickup_date ON reservations (user_id, pickup_date);