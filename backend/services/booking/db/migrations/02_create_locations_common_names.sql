CREATE TABLE locations_common_names (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    common_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Case-insensitive uniqueness per location
CREATE UNIQUE INDEX locations_common_names_location_name_unique
ON locations_common_names (location_id, lower(common_name));

-- Helps searching aliases
CREATE INDEX locations_common_names_location_id_idx
ON locations_common_names (location_id);