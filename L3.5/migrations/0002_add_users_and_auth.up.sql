-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    telegram_username VARCHAR(255) NULL,
    telegram_chat_id BIGINT NULL,
    preferred_notification VARCHAR(20) DEFAULT 'email',
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add user_id to bookings table
ALTER TABLE bookings ADD COLUMN user_id UUID NULL REFERENCES users(id) ON DELETE CASCADE;

-- Create index on user_id
CREATE INDEX idx_bookings_user_id ON bookings(user_id);

-- Create index on email for faster lookups
CREATE INDEX idx_users_email ON users(email);

-- Create index on telegram_username for bot registration
CREATE INDEX idx_users_telegram_username ON users(telegram_username) WHERE telegram_username IS NOT NULL;

-- Create index on role for faster queries
CREATE INDEX idx_users_role ON users(role);