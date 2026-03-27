-- migrations/000009_create_idempotency_keys.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS idempotency_keys ( 
    key                 VARCHAR(255)    PRIMARY KEY,
    request_hash        VARCHAR(64)     NOT NULL,
    response_status     INTEGER         NOT NULL,
    response_body       JSONB           NOT NULL,
    expires_at          TIMESTAMPTZ     NOT NULL,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Fast Lookup: Across All entries for a specific expire time 
CREATE INDEX IF NOT EXISTS idx_idempotency_key_expires_at
    ON idempotency_keys(expires_at);

COMMIT;