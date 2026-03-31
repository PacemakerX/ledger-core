-- migrations/000013_seed_platform_accounts.up.sql
BEGIN;

INSERT INTO accounts (id, account_number, customer_id, currency_id, type_id, country_id, is_active)
VALUES 
    ('99bf8606-446b-4603-8b9a-49a9e83003d4', 'ACC-PLATFORM-FLOAT',    NULL, 1, 3, 1, true),
    ('fa28695c-5204-49cf-84c1-06941abc9d4f', 'ACC-PLATFORM-CASH',     NULL, 1, 1, 1, true),
    ('5a20ea5e-feb5-45c5-9ec0-db460d2627bc', 'ACC-PLATFORM-REVENUE',  NULL, 1, 4, 1, true);

COMMIT;