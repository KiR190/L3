CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    channel SMALLINT NOT NULL,
    recipient TEXT NOT NULL,
    message TEXT NOT NULL,
    send_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0, 
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_send_at
    ON notifications (send_at);

CREATE INDEX IF NOT EXISTS idx_notifications_status
    ON notifications (status);