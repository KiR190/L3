-- Удаляем индексы
DROP INDEX IF EXISTS idx_items_occurred_at;
DROP INDEX IF EXISTS idx_items_category_id;
DROP INDEX IF EXISTS idx_items_type ;

-- Удаляем таблицу
DROP TABLE IF EXISTS items;