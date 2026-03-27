-- migrations/000005_create_customers.down.sql

BEGIN;

DROP TABLE IF EXISTS customers CASCADE;

COMMIT;