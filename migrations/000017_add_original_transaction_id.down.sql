DROP INDEX IF EXISTS idx_transactions_original_id;

ALTER TABLE transactions
DROP COLUMN IF EXISTS original_transaction_id;