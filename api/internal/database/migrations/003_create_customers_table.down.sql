-- Drop trigger
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;

-- Drop indexes
DROP INDEX IF EXISTS idx_customers_deleted_at;
DROP INDEX IF EXISTS idx_customers_email;

-- Drop table
DROP TABLE IF EXISTS customers;
