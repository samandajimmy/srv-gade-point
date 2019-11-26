ALTER TABLE vouchers
DROP COLUMN limit_per_user SMALLINT;

ALTER TABLE vouchers
DROP COLUMN day_purchase_limit SMALLINT;

ALTER TABLE vouchers
ADD COLUMN type SMALLINT;

