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