-- migrations/000012_create_indexes.down.sql

BEGIN;

DROP INDEX IF EXISTS idx_journal_entries_balance;

DROP INDEX IF EXISTS idx_exchange_rates_pair_date;

DROP INDEX IF EXISTS idx_accounts_active;

DROP INDEX IF EXISTS idx_transactions_pending;

DROP INDEX IF EXISTS idx_idempotency_expiry;

COMMIT;