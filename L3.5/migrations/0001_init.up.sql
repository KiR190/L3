-- Events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    event_date TIMESTAMP NOT NULL,
    total_seats INTEGER NOT NULL CHECK (total_seats > 0),
    payment_timeout_minutes INTEGER NOT NULL DEFAULT 30 CHECK (payment_timeout_minutes > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Bookings table
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    seats_count INTEGER NOT NULL CHECK (seats_count > 0),
    status VARCHAR(20) NOT NULL CHECK (status IN ('unpaid', 'paid', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    paid_at TIMESTAMP NULL,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Indexes for performance
CREATE INDEX idx_bookings_event_id ON bookings(event_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
CREATE INDEX idx_bookings_created_at ON bookings(created_at);
CREATE INDEX idx_events_event_date ON events(event_date);

-- Composite index for worker queries
CREATE INDEX idx_bookings_status_expires ON bookings(status, expires_at) WHERE status = 'unpaid';

