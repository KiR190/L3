-- Удаляем триггер
DROP TRIGGER IF EXISTS trg_item_audit ON items;

-- Удаляем триггерную функцию
DROP FUNCTION IF EXISTS fn_item_audit();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_item_history_role;
DROP INDEX IF EXISTS idx_item_history_created_at;
DROP INDEX IF EXISTS idx_item_history_item_id;

-- Удаляем таблицы в обратном порядке создания
DROP TABLE IF EXISTS item_history;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;