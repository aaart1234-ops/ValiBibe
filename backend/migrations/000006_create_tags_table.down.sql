-- удалить опциональные/вспомогательные индексы (безопасно даже если не существуют)
DROP INDEX IF EXISTS idx_tags_user_name_lower;
DROP INDEX IF EXISTS idx_tags_user_id;
DROP INDEX IF EXISTS ux_tags_user_name;

DROP TABLE IF EXISTS tags;
