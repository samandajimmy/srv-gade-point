ALTER TABLE vouchers
ADD COLUMN value DECIMAL NOT NULL,
DROP COLUMN day_purchase_limit SMALLINT;



DELETE FROM public.promo_codes
WHERE voucher_id=19;
DELETE FROM public.promo_codes
WHERE voucher_id=21;
DELETE FROM public.promo_codes
WHERE voucher_id=22;
DELETE FROM public.promo_codes
WHERE voucher_id=23;
DELETE FROM public.promo_codes
WHERE voucher_id=24;
DELETE FROM public.promo_codes
WHERE voucher_id=25;
DELETE FROM public.promo_codes
WHERE voucher_id=26;
DELETE FROM public.promo_codes
WHERE voucher_id=27;
DELETE FROM public.promo_codes
WHERE voucher_id=28;
DELETE FROM public.promo_codes
WHERE voucher_id=29;
DELETE FROM public.promo_codes
WHERE voucher_id=30;
DELETE FROM public.promo_codes
WHERE voucher_id=31;
DELETE FROM public.promo_codes
WHERE voucher_id=32;
