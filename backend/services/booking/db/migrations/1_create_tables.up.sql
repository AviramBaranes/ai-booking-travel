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