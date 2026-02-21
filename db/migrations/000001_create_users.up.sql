CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE subscription_status AS ENUM (
    'free',
    'active',
    'cancelled',
    'past_due'
);

CREATE TABLE users (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                   TEXT NOT NULL UNIQUE,
    password_hash           TEXT NOT NULL,
    name                    TEXT NOT NULL,
    role                    TEXT NOT NULL DEFAULT 'user'
                                CHECK (role IN ('user', 'admin')),
    stripe_customer_id      TEXT UNIQUE,
    stripe_subscription_id  TEXT UNIQUE,
    subscription_status     subscription_status NOT NULL DEFAULT 'free',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_stripe_customer_id ON users (stripe_customer_id);
CREATE INDEX idx_users_stripe_subscription_id ON users (stripe_subscription_id);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();
