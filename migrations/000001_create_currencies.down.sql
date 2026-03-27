-- migrations/000001_create_currencies.down.sql

BEGIN;

DROP TABLE IF EXISTS currencies CASCADE; 

COMMIT;