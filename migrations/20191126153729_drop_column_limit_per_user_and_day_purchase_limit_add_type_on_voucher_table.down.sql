ALTER TABLE vouchers
ADD COLUMN limit_per_user SMALLINT;

ALTER TABLE vouchers
ADD COLUMN day_purchase_limit SMALLINT;

ALTER TABLE vouchers
DROP COLUMN type;

