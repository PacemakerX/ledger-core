-- migrations/000006_create_accounts.down.sql

BEGIN;

DROP TABLE IF EXISTS accounts CASCADE;

DROP SEQUENCE IF EXISTS account_number_seq;

COMMIT;