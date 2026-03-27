-- migrations/000002_create_exchange_rates.down.sql

BEGIN;

DROP TABLE IF EXISTS exchange_rates;

COMMIT;