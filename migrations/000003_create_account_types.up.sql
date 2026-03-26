-- migrations/000003_create_account_types.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS account_types (
    id             SERIAL       PRIMARY KEY,
    name           VARCHAR(20)  UNIQUE NOT NULL,
    normal_balance VARCHAR(6)   NOT NULL,
    description    TEXT,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_account_type_name 
        CHECK (name IN ('asset','liability','equity','revenue','expense')),
    
    CONSTRAINT chk_normal_balance 
        CHECK (normal_balance IN ('DEBIT','CREDIT'))
);

-- Fast lookup: all rates for a specific account type
CREATE INDEX IF NOT EXISTS idx_account_types_name 
    ON account_types(name);

-- Seed account types
-- Assets and Expenses normally have DEBIT balance
-- Liabilities, Equity, Revenue normally have CREDIT balance
INSERT INTO account_types (name, normal_balance, description)
VALUES
    ('asset',     'DEBIT',  'Resources owned by the business'),
    ('liability', 'CREDIT', 'Obligations owed to external parties'),
    ('equity',    'CREDIT', 'Owner residual interest in assets'),
    ('revenue',   'CREDIT', 'Income earned from business operations'),
    ('expense',   'DEBIT',  'Costs incurred during business operations');

COMMIT;