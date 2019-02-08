CREATE TABLE IF NOT EXISTS article_categories (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    article_id INTEGER,
    category_id INTEGER,
    updated_at BIGINT DEFAULT NULL,
    created_at BIGINT DEFAULT NULL
);