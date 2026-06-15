-- 
-- Create location_aliases table to store common names for locations
--

CREATE TABLE location_aliases (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    alias TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX location_aliases_unique
ON location_aliases (location_id, lower(alias));

CREATE INDEX location_aliases_location_id_idx
ON location_aliases (location_id);

CREATE INDEX location_aliases_search_idx
ON location_aliases (lower(alias));

--
-- Add Location Is Airport Column default to true if iata not null else false
--
ALTER TABLE locations
ADD COLUMN is_airport BOOLEAN NOT NULL DEFAULT FALSE;

UPDATE locations
SET is_airport = TRUE
WHERE iata IS NOT NULL;