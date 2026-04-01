-- migrations/000014_fix_idempotency_keys.down.sql
BEGIN;

ALTER TABLE idempotency_keys 
    ALTER COLUMN response_status TYPE INTEGER USING response_status::integer,
    ALTER COLUMN response_body TYPE JSONB USING response_body::jsonb;

COMMIT;