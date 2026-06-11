CREATE TYPE location_alias_type AS ENUM ('translation', 'typo');

CREATE TABLE location_aliases (
    id BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type location_alias_type NOT NULL,
    language_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX location_aliases_unique
ON location_aliases (location_id, language_code, type, lower(name));

CREATE INDEX location_aliases_location_id_idx
ON location_aliases (location_id);

CREATE INDEX location_aliases_search_idx
ON location_aliases (lower(name));