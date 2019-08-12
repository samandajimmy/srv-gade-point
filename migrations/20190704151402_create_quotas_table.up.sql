-- Table: quotas

CREATE TABLE IF NOT EXISTS quotas (
    id SERIAL PRIMARY KEY NOT NULL,
    number_of_days INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    is_per_user SMALLINT NOT NULL DEFAULT 0,
    reward_id INTEGER REFERENCES rewards(id),
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_quotas ON quotas (id);
