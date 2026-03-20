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