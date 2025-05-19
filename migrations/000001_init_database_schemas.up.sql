CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(32) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TYPE frequency AS ENUM ('hourly', 'daily');

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP,

    city VARCHAR(32) NOT NULL,
    frequency VARCHAR(10) NOT NULL DEFAULT 'daily',

    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    confirm_token VARCHAR(64) UNIQUE,
    token_expires TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP
);