-- Удаляем индексы
DROP INDEX IF EXISTS idx_comments_text;
DROP INDEX IF EXISTS idx_comments_created_at;

-- Удаляем таблицу
DROP TABLE IF EXISTS comments;