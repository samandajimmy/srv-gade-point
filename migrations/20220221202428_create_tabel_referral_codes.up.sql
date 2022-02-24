CREATE TABLE IF NOT EXISTS referral_codes (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(10),
    referral_code VARCHAR(10),
    campaign_id SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_referral_codes
      FOREIGN KEY(campaign_id) 
	  REFERENCES campaigns(id)
);

CREATE INDEX index_referral_codes ON referral_codes (id, referral_code, campaign_id);
CREATE UNIQUE INDEX unique_index_referral_codes ON referral_codes (referral_code);