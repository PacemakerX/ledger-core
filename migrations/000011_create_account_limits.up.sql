-- migrations/000011_create_account_limits.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS accounts_limits( 
   id                UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
   account_id        UUID            NOT NULL REFERENCES accounts(id),
   limit_type        VARCHAR(20)     NOT NULL,
   max_amount        BIGINT          NOT NULL,
   current_usage     BIGINT          NOT NULL DEFAULT 0,
   reset_at          TIMESTAMPTZ,
   created_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
   updated_at        TIMESTAMPTZ     NOT NULL DEFAULT NOW()

   CONSTRAINT chk_limit_type 
        CHECK ( limit_type IN ('DAILY','MONTHLY','YEARLY')),
   CONSTRAINT chk_max_amount
        CHECK ( max_amount > 0),
   CONSTRAINT chk_current_usage
        CHECK ( current_usage >=0),
   CONSTRAINT uq_account_limit_type
        UNIQUE (account_id,limit_type)
);

-- Fast Lookup: All recorcds for a specific limit account id
CREATE INDEX IF NOT EXISTS idx_accounts_limit_accounts_id
    ON accounts_limits(account_id);

-- Fast Lookup: All recorcds for a specific limit reset time 
CREATE INDEX IF NOT EXISTS idx_accounts_limit_reset_at
    ON accounts_limits(reset_at);

COMMIT;