-- Extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Таблица items
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    quantity INT NOT NULL DEFAULT 0,
    location TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Таблица истории
CREATE TABLE IF NOT EXISTS item_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID,
    action TEXT NOT NULL, -- INSERT/UPDATE/DELETE
    old_data JSONB,       -- NULL for insert
    new_data JSONB,       -- NULL for delete
    user_id UUID,
    username TEXT,
    role TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_item_history_item_id ON item_history(item_id);
CREATE INDEX IF NOT EXISTS idx_item_history_created_at ON item_history(created_at);
CREATE INDEX IF NOT EXISTS idx_item_history_role ON item_history(role);

-- Триггерная функция, которая записывает историю
CREATE OR REPLACE FUNCTION fn_item_audit()
RETURNS TRIGGER AS $$
DECLARE
    v_user_id UUID := NULL;
    v_username TEXT := NULL;
    v_role TEXT := NULL;
BEGIN
    BEGIN
        v_user_id := current_setting('application.user_id', true)::UUID;
    EXCEPTION WHEN OTHERS THEN
        v_user_id := NULL;
    END;
    BEGIN
        v_username := current_setting('application.username', true);
    EXCEPTION WHEN OTHERS THEN
        v_username := NULL;
    END;
    BEGIN
        v_role := current_setting('application.role', true);
    EXCEPTION WHEN OTHERS THEN
        v_role := NULL;
    END;

    IF TG_OP = 'INSERT' THEN
        INSERT INTO item_history(item_id, action, old_data, new_data, user_id, username, role, created_at)
        VALUES (NEW.id, 'INSERT', NULL, to_jsonb(NEW), v_user_id, v_username, v_role, now());
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO item_history(item_id, action, old_data, new_data, user_id, username, role, created_at)
        VALUES (NEW.id, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW), v_user_id, v_username, v_role, now());
        NEW.updated_at := now();
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO item_history(item_id, action, old_data, new_data, user_id, username, role, created_at)
        VALUES (OLD.id, 'DELETE', to_jsonb(OLD), NULL, v_user_id, v_username, v_role, now());
        RETURN OLD;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_item_audit
    BEFORE INSERT OR UPDATE OR DELETE ON items
    FOR EACH ROW EXECUTE PROCEDURE fn_item_audit();
