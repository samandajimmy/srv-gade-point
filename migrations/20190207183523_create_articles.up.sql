CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(45),
    content TEXT,
    author_id INTEGER DEFAULT 0,
    updated_at BIGINT DEFAULT NULL,
    created_at BIGINT DEFAULT NULL
);