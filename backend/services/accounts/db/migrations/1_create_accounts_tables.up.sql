CREATE TYPE user_role AS ENUM ('customer', 'agent', 'admin', 'accountant');

CREATE TABLE
    organizations (
        id BIGSERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        is_organic BOOLEAN NOT NULL,
        icount_client_id INTEGER,
        phone TEXT,
        address TEXT,
        obligo INTEGER,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT organizations_icount_client_id_organic CHECK (
            (is_organic = TRUE AND icount_client_id IS NOT NULL)
            OR (is_organic = FALSE AND icount_client_id IS NULL)
        )
    );

CREATE TABLE
    offices (
        id BIGSERIAL PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
        organization_id BIGINT NOT NULL REFERENCES organizations (id) ON DELETE CASCADE,
        icount_client_id INTEGER,
        phone TEXT,
        address TEXT,
        obligo INTEGER,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE TABLE
    contacts (
        id BIGSERIAL PRIMARY KEY,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        role TEXT NOT NULL,
        cellphone TEXT NOT NULL,
        email TEXT NOT NULL,
        office_id BIGINT REFERENCES offices (id) ON DELETE CASCADE,
        organization_id BIGINT REFERENCES organizations (id) ON DELETE CASCADE,
        is_payment_responsible BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT contact_belongs_to_one CHECK (
            (
                office_id IS NOT NULL
                AND organization_id IS NULL
            )
            OR (
                office_id IS NULL
                AND organization_id IS NOT NULL
            )
        )
    );

CREATE TABLE
    users (
        id BIGSERIAL PRIMARY KEY,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        role user_role NOT NULL,
        email TEXT NOT NULL UNIQUE,
        phone_number TEXT UNIQUE,
        otp TEXT,
        office_id BIGINT REFERENCES offices (id),
        password_hash TEXT NOT NULL,
        last_login TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        CONSTRAINT users_office_agent_only CHECK (
                role = 'agent'
                OR office_id IS NULL
            )
    );

CREATE TABLE
    refresh_tokens (
        jti TEXT PRIMARY KEY,
        user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        admin_ref_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
        expires_at TIMESTAMPTZ NOT NULL
    );

-- ---------------------------------------------------------
-- Indexes
-- ---------------------------------------------------------

CREATE INDEX idx_offices_organization_id ON offices (organization_id);
CREATE INDEX idx_contacts_office_id ON contacts (office_id);
CREATE INDEX idx_contacts_organization_id ON contacts (organization_id);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE INDEX idx_users_office_id ON users (office_id) WHERE office_id IS NOT NULL;

CREATE INDEX idx_users_role_agent ON users (created_at DESC) WHERE role = 'agent';