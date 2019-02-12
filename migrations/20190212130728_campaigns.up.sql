
-- Type: status
CREATE TYPE IF NOT EXISTS status AS ENUM ('active', 'inactive');

-- Table: campaigns
CREATE TABLE IF NOT EXISTS campaigns (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status status DEFAULT 'inactive',
    validators JSONB NOT NULL, 
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);