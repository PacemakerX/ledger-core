-- migrations/000002_create_exchange_rates.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS exchange_rates (
    id              SERIAL          PRIMARY KEY,
    from_currency   INTEGER         NOT NULL REFERENCES currencies(id),
    to_currency     INTEGER         NOT NULL REFERENCES currencies(id),
    rate            NUMERIC(20,8)   NOT NULL,
    source          VARCHAR(10)     NOT NULL,
    effective_date  DATE            NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_source_of_rates
        CHECK (source IN ('manual', 'api')),

    CONSTRAINT chk_rate
        CHECK (rate > 0),

    CONSTRAINT chk_different_currencies
        CHECK (from_currency != to_currency),

    CONSTRAINT uq_exchange_rate_per_day
        UNIQUE (from_currency, to_currency, effective_date)
);

-- Fast lookup: "give me USD to INR rate for today"
CREATE INDEX IF NOT EXISTS idx_exchange_rates_lookup
    ON exchange_rates(from_currency, to_currency, effective_date DESC);

-- Fast lookup: all rates for a specific date
CREATE INDEX IF NOT EXISTS idx_exchange_rates_date
    ON exchange_rates(effective_date DESC);

COMMIT;