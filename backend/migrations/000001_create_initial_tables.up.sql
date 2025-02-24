CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       nickname TEXT NOT NULL,
                       email TEXT UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       subscription_status TEXT NULL
);

CREATE TABLE notes (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       user_id UUID NOT NULL,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       archived_at TIMESTAMP NULL,
                       memory_level INT DEFAULT 0,
                       deleted_at TIMESTAMP NULL,
                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE folders (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         user_id UUID NOT NULL,
                         name TEXT NOT NULL,
                         parent_id UUID NULL,
                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                         FOREIGN KEY (parent_id) REFERENCES folders(id) ON DELETE SET NULL
);

CREATE TABLE tags (
                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      user_id UUID NOT NULL,
                      name TEXT NOT NULL UNIQUE,
                      FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE note_folders (
                              note_id UUID NOT NULL,
                              folder_id UUID NOT NULL,
                              PRIMARY KEY (note_id, folder_id),
                              FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
                              FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE TABLE note_tags (
                           note_id UUID NOT NULL,
                           tag_id UUID NOT NULL,
                           PRIMARY KEY (note_id, tag_id),
                           FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
                           FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE TABLE reminders (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           user_id UUID NOT NULL,
                           schedule INTERVAL NOT NULL,
                           next_trigger TIMESTAMP NOT NULL,
                           FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE card_reviews (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              note_id UUID NOT NULL,
                              user_id UUID NOT NULL,
                              reviewed_at TIMESTAMP DEFAULT NOW(),
                              remembered BOOLEAN NOT NULL,
                              FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
