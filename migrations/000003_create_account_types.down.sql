-- migrations/000003_create_account_types.down.sql

BEGIN;

DROP TABLE  IF EXISTS account_types CASCADE;

COMMIT;
