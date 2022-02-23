CREATE TABLE IF NOT EXISTS referral_codes (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(50),
    referral_code VARCHAR(20),
    campaign_id SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_referral_codes ON referral_codes (id, referral_code, campaign_id);