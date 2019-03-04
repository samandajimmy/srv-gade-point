ALTER TABLE campaign_transactions
DROP COLUMN promo_code_id,
ADD COLUMN voucher_id SMALLINT REFERENCES vouchers(id) NULL;