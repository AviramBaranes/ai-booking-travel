CREATE TYPE user_role AS ENUM (
    'customer',
    'agent',
    'admin'
);

CREATE TABLE users (
    id             SERIAL                 PRIMARY KEY,
    role           user_role              NOT NULL,
    username       VARCHAR(255)        NOT NULL UNIQUE,
    phone_number   VARCHAR(20)            UNIQUE,
    otp            VARCHAR(6)             UNIQUE,
    office_code    VARCHAR(10),
    agent_code     VARCHAR(10),
    password_hash  VARCHAR(255)           NOT NULL,
    last_login     TIMESTAMP WITH TIME ZONE,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    jti         TEXT           PRIMARY KEY,
    user_id     SERIAL         NOT NULL REFERENCES users(id),
    expires_at  TIMESTAMPTZ    NOT NULL
);
