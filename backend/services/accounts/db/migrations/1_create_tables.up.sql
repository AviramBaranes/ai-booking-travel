CREATE TYPE user_role AS ENUM ('customer', 'agent', 'admin');

CREATE TABLE
    organizations (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        is_organic BOOLEAN NOT NULL,
        phone VARCHAR(20),
        address TEXT,
        obligo DECIMAL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    offices (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        organization_id INTEGER NOT NULL REFERENCES organizations (id),
        phone VARCHAR(20),
        address TEXT,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    contacts (
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
        role TEXT NOT NULL,
        cellphone VARCHAR(20) NOT NULL,
        email VARCHAR(255) NOT NULL,
        office_id INTEGER REFERENCES offices (id),
        organization_id INTEGER REFERENCES organizations (id),
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
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
        id SERIAL PRIMARY KEY,
        role user_role NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        phone_number VARCHAR(20) UNIQUE,
        otp VARCHAR(6) UNIQUE,
        office_id INTEGER REFERENCES offices (id),
        password_hash VARCHAR(255) NOT NULL,
        last_login TIMESTAMP
        WITH
            TIME ZONE,
            created_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP
        WITH
            TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT users_office_agent_only CHECK (
                role = 'agent'
                OR office_id IS NULL
            )
    );

CREATE TABLE
    refresh_tokens (
        jti TEXT PRIMARY KEY,
        user_id SERIAL NOT NULL REFERENCES users (id),
        admin_ref_id INTEGER REFERENCES users (id),
        expires_at TIMESTAMPTZ NOT NULL
    );