DROP INDEX IF EXISTS idx_notes_folder_id;
ALTER TABLE notes DROP CONSTRAINT IF EXISTS fk_notes_folder;
ALTER TABLE notes DROP COLUMN IF EXISTS folder_id;
