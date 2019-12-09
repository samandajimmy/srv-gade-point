-- Table: referral_transactions

CREATE TABLE IF NOT EXISTS referral_transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(50),
    ref_id VARCHAR(50),
    used_referral_code VARCHAR(20),
    type SMALLINT NOT NULL DEFAULT 0,
    reward_referral INTEGER,
    reward_type VARCHAR(255),
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_referral_transactions ON referral_transactions (id, used_referral_code);