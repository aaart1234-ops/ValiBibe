-- Удаляем default-значения у временных полей
ALTER TABLE users
    ALTER COLUMN created_at DROP DEFAULT;

ALTER TABLE notes
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN updated_at DROP DEFAULT;

-- Обновляем nullability: делаем поля NOT NULL (если меняем с *time.Time на time.Time)
ALTER TABLE notes
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;