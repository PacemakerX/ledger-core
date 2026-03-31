-- migrations/000006_create_accounts.up.sql

BEGIN;

CREATE SEQUENCE IF NOT EXISTS account_number_seq;

CREATE TABLE IF NOT EXISTS accounts(

    id                      UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    account_number          VARCHAR(20)     UNIQUE NOT NULL DEFAULT 'ACC-' || LPAD(nextval('account_number_seq')::text, 6, '0'),
    customer_id             UUID            REFERENCES customers(id),
    currency_id             INTEGER         NOT NULL REFERENCES currencies(id),
    type_id                 INTEGER         NOT NULL REFERENCES account_types(id),
    country_id              INTEGER         NOT NULL REFERENCES countries(id),
    is_active               BOOLEAN         NOT NULL DEFAULT true,
    daily_debit_limit       BIGINT          DEFAULT NULL,
    daily_credit_limit      BIGINT          DEFAULT NULL,
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_customer_account
        UNIQUE (customer_id, type_id, currency_id)
);

-- Fast lookup: all rates for a specific customer id
CREATE INDEX IF NOT EXISTS idx_accounts_customer_id
    ON accounts(customer_id);

-- Fast lookup: all rates for a specific currency
CREATE INDEX IF NOT EXISTS idx_accounts_currency_id
    on accounts(currency_id);

-- Fast lookup: all rates for a specific account type
CREATE INDEX IF NOT EXISTS idx_accounts_type_id
    on accounts(type_id);

-- Fast lookup: all rates for a specific account number
CREATE INDEX IF NOT EXISTS idx_accounts_account_number
    on accounts(account_number);

COMMIT;