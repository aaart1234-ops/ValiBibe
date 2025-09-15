-- folders
CREATE TABLE IF NOT EXISTS folders (
                                       id UUID PRIMARY KEY,
                                       user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    parent_id UUID NULL REFERENCES folders(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- имя уникально в пределах (user_id, parent_id)
CREATE UNIQUE INDEX IF NOT EXISTS ux_folders_user_parent_name
    ON folders (user_id, parent_id, name);

-- индексы под дерево/поиск
CREATE INDEX IF NOT EXISTS idx_folders_user_parent ON folders (user_id, parent_id);
CREATE INDEX IF NOT EXISTS idx_folders_parent ON folders (parent_id);
CREATE INDEX IF NOT EXISTS idx_folders_user_name ON folders (user_id, name);

-- триггер для updated_at (опционально, если не обновляете в коде)
CREATE OR REPLACE FUNCTION set_folders_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_folders_updated_at
    BEFORE UPDATE ON folders
    FOR EACH ROW
    EXECUTE FUNCTION set_folders_updated_at();
