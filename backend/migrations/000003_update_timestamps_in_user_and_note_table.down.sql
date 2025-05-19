-- Возвращаем default-значения
ALTER TABLE users
    ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE notes
    ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

-- Делаем поля nullable (если возвращаем к *time.Time)
ALTER TABLE notes
    ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;
