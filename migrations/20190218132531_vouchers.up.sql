-- Table: vouchers
-- for status 0 --> INACTIVE and 1 --> ACTIVE
CREATE TABLE IF NOT EXISTS vouchers (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    point SMALLINT NOT NULL,
    journal_account CHAR(20) NOT NULL,
    value DECIMAL NOT NULL,
    image_url VARCHAR(255) NOT NULL, 
    status SMALLINT DEFAULT 0,
    stock SMALLINT NOT NULL,
    prefix_promo_code CHAR(5) NOT NULL,
    validators JSONB NOT NULL, 
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_vouchers ON vouchers (name, start_date, end_date, status);
