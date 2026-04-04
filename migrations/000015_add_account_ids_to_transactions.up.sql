ALTER TABLE transactions
ADD COLUMN from_account_id UUID REFERENCES accounts(id),
ADD COLUMN to_account_id UUID REFERENCES accounts(id);

CREATE INDEX idx_transactions_from_account_id ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account_id ON transactions(to_account_id);