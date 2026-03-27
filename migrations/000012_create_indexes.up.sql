-- migrations/000012_create_indexes.up.sql

-- migrations/000012_create_indexes.up.sql
BEGIN;

-- Balance calculation query optimization
-- SELECT SUM(amount) FROM journal_entries 
-- WHERE account_id = $1 AND entry_type = 'DEBIT'
CREATE INDEX IF NOT EXISTS idx_journal_entries_balance
    ON journal_entries(account_id, entry_type);

-- Exchange rate lookup optimization
-- "Give me latest USD to INR rate"
CREATE INDEX IF NOT EXISTS idx_exchange_rates_pair_date
    ON exchange_rates(from_currency, to_currency, effective_date DESC);

-- Active accounts lookup
CREATE INDEX IF NOT EXISTS idx_accounts_active
    ON accounts(is_active) WHERE is_active = true;

-- Pending transactions cleanup
CREATE INDEX IF NOT EXISTS idx_transactions_pending
    ON transactions(status, created_at) WHERE status = 'PENDING';

-- Expired idempotency keys cleanup
CREATE INDEX IF NOT EXISTS idx_idempotency_expiry
    ON idempotency_keys(expires_at);

COMMIT;