CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
    author VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- индекс для быстрого поиска по тексту
CREATE INDEX IF NOT EXISTS idx_comments_text ON comments USING gin(to_tsvector('russian', text));

-- индекс для сортировки
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);
