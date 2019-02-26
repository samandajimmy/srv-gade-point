-- Table: promo_codes
/* 0 -> available
    1 -> redeemed
    2 -> used
    3 -> expired*/

CREATE TABLE IF NOT EXISTS promo_codes (
    promo_code VARCHAR(10) PRIMARY KEY NOT NULL,
    status SMALLINT DEFAULT 0,
    user_id VARCHAR(50),
    voucher_id SMALLINT REFERENCES vouchers(id) NOT NULL,
    redeemed_date TIMESTAMP DEFAULT NULL,
    used_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_promo_codes ON promo_codes (promo_code, status, user_id);
