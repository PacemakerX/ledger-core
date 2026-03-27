-- migrations/000006_create_accounts.down.sql

BEGIN;

DROP TABLE IF EXISTS accounts;

DROP SEQUENCE IF EXISTS account_number_seq;

COMMIT;