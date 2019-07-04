-- Table: tags

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_tags ON tags (id);

-- Table: reward_tags

CREATE TABLE IF NOT EXISTS reward_tags (
    reward_id INTEGER REFERENCES rewards (id) ON UPDATE CASCADE ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags (id) ON UPDATE CASCADE,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT reward_tags_pkey PRIMARY KEY (reward_id, tag_id)
);

CREATE INDEX index_reward_tags ON reward_tags (reward_id, tag_id);
