-- migrations/000010_create_audit_logs.down.sql

BEGIN;

DROP TABLE IF EXISTS audit_logs CASCADE;

COMMIT;