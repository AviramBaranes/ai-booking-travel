-- Postgres will not let a newly added enum value be used in the transaction that adds it, so
-- nothing below may refer to 'hertz_gr'.
ALTER TYPE broker ADD VALUE IF NOT EXISTS 'hertz_gr';

-- brands lists the supplier brands a station serves, as OTA vendor codes (e.g. {ZE,ZT}), for a
-- broker that sells several brands. It is empty for brokers without brands.
ALTER TABLE location_broker_codes ADD COLUMN brands TEXT[] NOT NULL DEFAULT '{}';

-- A broker without brands keeps one station per location. A broker with brands can hold one per
-- brand, since a brand may work from its own desk at the same airport: Hertz Greece serves
-- Mykonos Airport from GR001 for Hertz and from GR221 for Thrifty and Dollar. The import keeps
-- the brands of one location's stations apart.
DROP INDEX uniq_location_broker;
CREATE UNIQUE INDEX uniq_location_broker ON location_broker_codes (location_id, broker)
WHERE cardinality(brands) = 0;
