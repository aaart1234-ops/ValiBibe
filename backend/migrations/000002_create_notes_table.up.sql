CREATE TABLE notes (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       memory_level INT DEFAULT 0,
                       archived BOOLEAN DEFAULT FALSE,
                       next_review_at TIMESTAMP NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индекс для ускорения выборки всех заметок пользователя
CREATE INDEX idx_notes_user_id ON notes(user_id);

-- Индекс для ускорения выборки карточек по memory_level (например, при тестировании)
CREATE INDEX idx_notes_memory_level ON notes(memory_level);