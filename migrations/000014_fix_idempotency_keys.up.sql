-- migrations/000014_fix_idempotency_keys.up.sql
BEGIN;

ALTER TABLE idempotency_keys 
    ALTER COLUMN response_status TYPE VARCHAR(20),
    ALTER COLUMN response_body TYPE TEXT;

COMMIT;