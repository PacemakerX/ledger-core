-- migrations/000007_create_transactions.down.sql

BEGIN;

DROP TABLE IF EXISTS transactions CASCADE;

COMMIT;