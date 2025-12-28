-- Drop indexes
DROP INDEX IF EXISTS idx_bookings_status_expires;
DROP INDEX IF EXISTS idx_events_event_date;
DROP INDEX IF EXISTS idx_bookings_created_at;
DROP INDEX IF EXISTS idx_bookings_expires_at;
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_event_id;

-- Drop tables
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS events;

