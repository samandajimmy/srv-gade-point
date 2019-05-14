ALTER TABLE campaign_transactions DROP CONSTRAINT campaign_transactions_promo_code_id_fkey;

ALTER TABLE campaign_transactions
DROP COLUMN promo_code_id;

DROP TABLE IF EXISTS promo_codes;

-- Table: voucher_codes
/*  Generated  by admin                 = 0 -> available
    user bought a voucher with points   = 1 -> bought
    the voucher has been used           = 2 -> redeemed
    voucher expired                     = 3 -> expired*/

CREATE TABLE IF NOT EXISTS voucher_codes (
    id SERIAL PRIMARY KEY NOT NULL,
    promo_code VARCHAR(10) UNIQUE,
    status SMALLINT DEFAULT 0,
    user_id VARCHAR(50),
    voucher_id INT REFERENCES vouchers(id) NOT NULL,
    redeemed_date TIMESTAMP DEFAULT NULL,
    bought_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE campaign_transactions
ADD COLUMN voucher_code_id INT REFERENCES voucher_codes(id) NULL;

CREATE INDEX index_voucher_codes ON voucher_codes (promo_code, status, user_id);