-- migrations/000009_create_idempotency_keys.down.sql

BEGIN;

DROP TABLE IF EXISTS idempotency_keys CASCADE;

COMMIT;