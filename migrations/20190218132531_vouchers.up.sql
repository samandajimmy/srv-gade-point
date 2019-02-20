-- Table: vouchers
-- for status 0 --> INACTIVE and 1 --> ACTIVE
CREATE TABLE IF NOT EXISTS vouchers (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    point INTEGER NOT NULL,
    journal_account VARCHAR(255) NOT NULL,
    value DECIMAL NOT NULL,
    image_url VARCHAR(255) NOT NULL, 
    status SMALLINT DEFAULT 0,
    validators JSONB NOT NULL, 
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_vouchers ON vouchers (id, name, start_date, end_date, status);
