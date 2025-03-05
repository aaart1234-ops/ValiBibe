-- Удаление индексов перед удалением таблиц
DROP INDEX IF EXISTS idx_notes_user_id;
DROP INDEX IF EXISTS idx_note_tags_tag_id;
DROP INDEX IF EXISTS idx_note_tags_note_id;
DROP INDEX IF EXISTS idx_note_folders_folder_id;
DROP INDEX IF EXISTS idx_note_folders_note_id;
DROP INDEX IF EXISTS idx_card_reviews_note_id;
DROP INDEX IF EXISTS idx_card_reviews_user_id;
DROP INDEX IF EXISTS idx_reminders_user_id;
DROP INDEX IF EXISTS idx_tags_user_id;
DROP INDEX IF EXISTS idx_folders_user_id;

DROP INDEX IF EXISTS idx_notes_memory_level;
DROP INDEX IF EXISTS idx_card_reviews_reviewed_at;
DROP INDEX IF EXISTS idx_folders_parent_id;
