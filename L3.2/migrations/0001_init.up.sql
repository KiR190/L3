CREATE TABLE IF NOT EXISTS short_urls (
    id           SERIAL PRIMARY KEY,
    short_code   VARCHAR(16) UNIQUE NOT NULL,
    original    TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS click_events (
    id           SERIAL PRIMARY KEY,
    short_url_id INT NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    user_agent  TEXT,
    timestamp   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_click_events_short_id ON click_events(short_url_id);
