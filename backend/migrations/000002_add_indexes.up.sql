-- Индексы для ускорения работы с внешними ключами
CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_note_tags_tag_id ON note_tags(tag_id);
CREATE INDEX idx_note_tags_note_id ON note_tags(note_id);
CREATE INDEX idx_note_folders_folder_id ON note_folders(folder_id);
CREATE INDEX idx_note_folders_note_id ON note_folders(note_id);
CREATE INDEX idx_card_reviews_note_id ON card_reviews(note_id);
CREATE INDEX idx_card_reviews_user_id ON card_reviews(user_id);
CREATE INDEX idx_reminders_user_id ON reminders(user_id);
CREATE INDEX idx_tags_user_id ON tags(user_id);
CREATE INDEX idx_folders_user_id ON folders(user_id);

-- Индексы для оптимизации выборок
CREATE INDEX idx_notes_memory_level ON notes(memory_level);
CREATE INDEX idx_card_reviews_reviewed_at ON card_reviews(reviewed_at);

-- Индекс для ускорения работы с деревом папок
CREATE INDEX idx_folders_parent_id ON folders(parent_id);
