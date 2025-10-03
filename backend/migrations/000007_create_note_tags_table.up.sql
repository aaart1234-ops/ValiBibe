CREATE TABLE IF NOT EXISTS note_tags (
   note_id UUID NOT NULL,
   tag_id  UUID NOT NULL,
   PRIMARY KEY (note_id, tag_id),
   CONSTRAINT fk_note_tags_note FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    CONSTRAINT fk_note_tags_tag  FOREIGN KEY (tag_id)  REFERENCES tags(id)  ON DELETE CASCADE
    );

-- индекс чтобы быстро находить все заметки по tag_id
CREATE INDEX IF NOT EXISTS idx_note_tags_tag_id ON note_tags (tag_id);
