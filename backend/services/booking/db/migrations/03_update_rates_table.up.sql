DELETE FROM hertz_markup_rates;

ALTER TABLE hertz_markup_rates RENAME TO markup_rates;

ALTER TABLE markup_rates RENAME COLUMN country TO country_code;

ALTER TABLE markup_rates DROP COLUMN brand;
ALTER TABLE markup_rates DROP COLUMN pickup_date_from;
ALTER TABLE markup_rates DROP COLUMN pickup_date_to;
ALTER TABLE markup_rates DROP COLUMN car_group;
ALTER TABLE markup_rates DROP COLUMN num_of_rental_days_from;
ALTER TABLE markup_rates DROP COLUMN num_of_rental_days_to;

ALTER TABLE markup_rates ADD COLUMN broker broker NOT NULL;

ALTER TABLE markup_rates
ADD CONSTRAINT markup_rates_country_code_broker_unique
UNIQUE (country_code, broker);