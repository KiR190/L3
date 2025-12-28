-- Drop indexes
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_telegram_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_bookings_user_id;

-- Remove user_id column from bookings
ALTER TABLE bookings DROP COLUMN IF EXISTS user_id;

-- Drop users table
DROP TABLE IF EXISTS users;

