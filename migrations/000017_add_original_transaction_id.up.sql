ALTER TABLE transactions
ADD COLUMN original_transaction_id UUID REFERENCES transactions(id);

CREATE INDEX idx_transactions_original_id ON transactions(original_transaction_id);