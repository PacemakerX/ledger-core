-- migrations/000011_create_account_limits.down.sql

BEGIN;

DROP TABLE IF EXISTS account_limits CASCADE;

COMMIT;