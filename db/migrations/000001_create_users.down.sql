DROP TRIGGER IF EXISTS set_users_updated_at ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS subscription_status;
DROP EXTENSION IF EXISTS "pgcrypto";
