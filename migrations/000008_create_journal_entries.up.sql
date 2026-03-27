-- migrations/000008_create_journal_entries.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS journal_entries ( 

    id                  UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id      UUID             NOT NULL REFERENCES transactions(id),
    account_id          UUID             NOT NULL REFERENCES accounts(id),
    entry_type          VARCHAR(6)       NOT NULL,
    amount              BIGINT           NOT NULL, 
    currency_id         INTEGER          NOT NULL REFERENCES currencies(id),
    exchange_rate_id    INTEGER          REFERENCES exchange_rates(id),
    description         TEXT,
    created_at          TIMESTAMPTZ      NOT NULL DEFAULT NOW(),


    CONSTRAINT chk_entry_type
        CHECK (entry_type IN ('DEBIT','CREDIT')),
    CONSTRAINT chk_amount
        CHECK (amount > 0) 

);

-- Fast Lookup: All journal entries for a specific account id
CREATE INDEX IF NOT EXISTS idx_journal_entries_account_id
    ON journal_entries(account_id);

-- Fast Lookup: All journal entries for a speicific transaction id
CREATE INDEX IF NOT EXISTS idx_journal_entries_transaction_id
    ON journal_entries(transaction_id);

-- Fast Lookup: All journeal entries for a specific entry at a time
CREATE INDEX IF NOT EXISTS idx_journal_entries_created_at
    ON journal_entries(created_at DESC);

COMMIT;
