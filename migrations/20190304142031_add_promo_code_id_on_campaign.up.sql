ALTER TABLE campaign_transactions
ADD COLUMN promo_code_id SMALLINT REFERENCES promo_codes(id) NULL,
DROP COLUMN voucher_id;
