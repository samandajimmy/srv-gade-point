-- Table: promo_codes
/*  Generated  by admin                 = 0 -> available
    user bought a voucher with points   = 1 -> bought
    the voucher has been used           = 2 -> redeemed
    voucher expired                     = 3 -> expired*/

CREATE TABLE IF NOT EXISTS promo_codes (
    id SERIAL PRIMARY KEY NOT NULL,
    promo_code VARCHAR(10) UNIQUE,
    status SMALLINT DEFAULT 0,
    user_id VARCHAR(50),
    voucher_id SMALLINT REFERENCES vouchers(id) NOT NULL,
    redeemed_date TIMESTAMP DEFAULT NULL,
    bought_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_promo_codes ON promo_codes (id, promo_code, status, user_id);
