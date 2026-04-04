ALTER TABLE transactions
ADD COLUMN amount BIGINT,
ADD COLUMN currency_id INTEGER REFERENCES currencies(id);