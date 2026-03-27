-- migrations/000008_create_journal_entries.down.sql

BEGIN;

DROP TABLE IF EXISTS journal_entries CASCADE;

COMMIT; 