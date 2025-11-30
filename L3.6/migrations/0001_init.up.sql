CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',

    category_id TEXT NULL,
    note TEXT NULL,

    occurred_at TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс по дате
CREATE INDEX idx_items_occurred_at ON items(occurred_at);

-- Индекс по категории
CREATE INDEX idx_items_category_id ON items(category_id);

-- Индекс по типу
CREATE INDEX idx_items_type ON items(type);