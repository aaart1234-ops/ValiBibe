-- Добавляем уникальный индекс для тегов:
-- уникальность имени в пределах пользователя (case-insensitive)
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_userid_name_lower
    ON tags (user_id, LOWER(name));