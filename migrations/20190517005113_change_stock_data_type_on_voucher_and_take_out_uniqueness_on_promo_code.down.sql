ALTER TABLE vouchers
ALTER COLUMN stock TYPE SMALLINT;

ALTER TABLE voucher_codes
ADD CONSTRAINT voucher_codes_promo_code_key UNIQUE (promo_code);