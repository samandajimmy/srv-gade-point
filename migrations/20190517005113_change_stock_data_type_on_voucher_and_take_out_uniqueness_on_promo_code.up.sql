ALTER TABLE vouchers
ALTER COLUMN stock TYPE INTEGER;

ALTER TABLE voucher_codes
DROP CONSTRAINT voucher_codes_promo_code_key;