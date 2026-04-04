DROP INDEX IF EXISTS idx_transactions_from_account_id;
DROP INDEX IF EXISTS idx_transactions_to_account_id;

ALTER TABLE transactions
DROP COLUMN IF EXISTS from_account_id,
DROP COLUMN IF EXISTS to_account_id;