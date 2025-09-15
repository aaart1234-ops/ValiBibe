-- колонка folder_id у заметок
ALTER TABLE notes
    ADD COLUMN IF NOT EXISTS folder_id UUID NULL;

ALTER TABLE notes
    ADD CONSTRAINT fk_notes_folder
        FOREIGN KEY (folder_id) REFERENCES folders(id)
            ON DELETE CASCADE;

-- индекс для выборок по папке
CREATE INDEX IF NOT EXISTS idx_notes_folder_id ON notes(folder_id);
