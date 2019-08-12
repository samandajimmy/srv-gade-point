-- Table: rewards
-- for is_promo_code 0 --> without promo_code
--                   1 --> with promo_code
--
-- for type          0 --> point
--                   1 --> discount
--                   2 --> goldback
--                   3 --> voucher

CREATE TABLE IF NOT EXISTS rewards (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    terms_and_conditions TEXT,
    how_to_use TEXT,
    journal_account VARCHAR(20) NOT NULL,
    promo_code VARCHAR(20) NULL,
    is_promo_code SMALLINT NOT NULL DEFAULT 0,
    custom_period VARCHAR,
    type SMALLINT NOT NULL DEFAULT 0,
    validators JSONB NOT NULL,
    campaign_id INTEGER REFERENCES campaigns(id),
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_rewards ON rewards (name, journal_account, promo_code);
