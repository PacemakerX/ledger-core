-- migrations/000004_create_countries.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS countries(
    id              SERIAL          PRIMARY KEY,
    name            VARCHAR(50)     UNIQUE NOT NULL,
    iso_code        VARCHAR(30)     UNIQUE NOT NULL,
    dial_code       VARCHAR(5)      NOT NULL,
    currency_id     INTEGER         REFERENCES currencies(id),
    is_active       BOOLEAN         NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

INSERT INTO countries (id, name, iso_code, dial_code, currency_id)
VALUES (
    1,
    'India',
    'IN',
    '+91',
    (SELECT id FROM currencies WHERE code = 'INR')
);

COMMIT;