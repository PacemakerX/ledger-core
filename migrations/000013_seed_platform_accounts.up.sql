-- migrations/000013_seed_platform_accounts.up.sql
BEGIN;

INSERT INTO accounts (id, account_number, customer_id, currency_id, type_id, country_id, is_active)
VALUES 
    ('99bf8606-446b-4603-8b9a-49a9e83003d4', 'ACC-PLATFORM-FLOAT',    NULL, 1, 3, 1, true),
    ('fa28695c-5204-49cf-84c1-06941abc9d4f', 'ACC-PLATFORM-CASH',     NULL, 1, 1, 1, true),
    ('5a20ea5e-feb5-45c5-9ec0-db460d2627bc', 'ACC-PLATFORM-REVENUE',  NULL, 1, 4, 1, true);


-- Test customers
INSERT INTO customers (id, first_name, last_name, aadhar_number, country_id, phone_number, email, kyc_status, is_active)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Sparsh', 'Soni', '123412341234', 1, '9999999999', 'sparsh@test.com', 'verified', true),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Test',   'User', '432143214321', 1, '8888888888', 'test@test.com',   'verified', true);

-- Customer accounts
INSERT INTO accounts (id, account_number, customer_id, currency_id, type_id, country_id, is_active)
VALUES
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'ACC-001001', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, 1, 1, true),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'ACC-001002', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 1, 1, 1, true);

-- Seed initial balance for Sparsh (100000 paise = 1000 INR)
INSERT INTO transactions (id, idempotency_key, type, status, initiated_by)
VALUES ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'seed-balance-001', 'ADJUSTMENT', 'COMPLETED', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa');

INSERT INTO journal_entries (transaction_id, account_id, entry_type, amount, currency_id)
VALUES
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '99bf8606-446b-4603-8b9a-49a9e83003d4', 'CREDIT', 100000, 1),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'DEBIT',  100000, 1);

COMMIT;