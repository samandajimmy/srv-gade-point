ALTER TABLE vouchers
ADD COLUMN day_purchase_limit SMALLINT,
DROP COLUMN value;