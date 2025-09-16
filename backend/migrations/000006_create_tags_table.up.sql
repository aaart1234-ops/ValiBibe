CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, name) -- у одного пользователя имя тега уникальное
);

-- уникальность имени тега в рамках пользователя
CREATE UNIQUE INDEX IF NOT EXISTS ux_tags_user_name
    ON tags (user_id, name);

-- индекс для быстрого списка/фильтрации тегов пользователя
CREATE INDEX IF NOT EXISTS idx_tags_user_id ON tags (user_id);

-- опционально: индекс для case-insensitive поиска (ускоряет ILIKE / lower(name) = ...)
CREATE INDEX IF NOT EXISTS idx_tags_user_name_lower ON tags (user_id, lower(name));
