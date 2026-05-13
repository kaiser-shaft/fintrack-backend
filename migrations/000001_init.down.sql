BEGIN;

DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP INDEX IF EXISTS idx_accounts_user_id;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS category_type;

COMMIT;
