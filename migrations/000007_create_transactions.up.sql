-- migrations/000007_create_transactions.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS transactions ( 
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key     VARCHAR(255)    UNIQUE NOT NULL,
    type                VARCHAR(20)     NOT NULL,
    status              VARCHAR(20)     NOT NULL DEFAULT 'PENDING',
    initiated_by        UUID            NOT NULL,
    metadata            JSONB,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    completed_at        TIMESTAMPTZ,
    
    CONSTRAINT chk_transaction_type 
        CHECK( type IN ('TRANSFER','REFUND','ADJUSTMENT')),
    
    CONSTRAINT chk_transaction_status
        CHECK(status IN ('PENDING','COMPLETED','FAILED'))
);

-- Fast lookup: All transactions for a specific status 
CREATE INDEX IF NOT EXISTS idx_transaction_status
    ON transactions(status);

-- Fast lookup: All transactions for a specific idempotency key
CREATE INDEX IF NOT EXISTS idx_transaction_idempotency_key
    ON transactions(idempotency_key);

-- Fast lookup: All transactions for a specific initiator
CREATE INDEX IF NOT EXISTS idx_transaction_initiated_by
    ON transactions(initiated_by);

COMMIT;