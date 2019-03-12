-- Table: account_tokens
-- for status 0 --> INACTIVE and 1 --> ACTIVE
CREATE TABLE IF NOT EXISTS account_tokens (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    expire_at TIMESTAMP DEFAULT NULL,
    status SMALLINT DEFAULT 0, 
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_account_tokens ON account_tokens (username, status);
