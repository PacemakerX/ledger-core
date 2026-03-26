-- migrations/000005_create_customers.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS customers (
    id              UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name      VARCHAR(50)     NOT NULL,
    middle_name     VARCHAR(50),
    last_name       VARCHAR(50)     NOT NULL,
    aadhar_number   VARCHAR(12)     UNIQUE,
    country_id      INTEGER         NOT NULL REFERENCES countries(id),
    phone_number    VARCHAR(15)     NOT NULL,
    email           VARCHAR(255)    UNIQUE NOT NULL,
    verified        BOOLEAN         NOT NULL DEFAULT false,
    kyc_status      VARCHAR(20)     NOT NULL DEFAULT 'unverified',
    is_active       BOOLEAN NOT     NULL DEFAULT true,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()

    CONSTRAINT chk_ky_status 
        CHECK (kyc_status IN ('unverified','pending','verified','rejected'))
);

-- Fast lookup: all rates for a specific email id 
CREATE INDEX IF NOT EXISTS idx_customer_email ON customers(email);

-- Fast lookup: all rates for a specific phone number
CREATE INDEX IF NOT EXISTS idx_customer_phone ON customers(phone_number);

-- Fast lookup: all rates for a specific kyc_status
CREATE INDEX IF NOT EXISTS idx_kyc_status     ON customers(kyc_status);

COMMIT;
