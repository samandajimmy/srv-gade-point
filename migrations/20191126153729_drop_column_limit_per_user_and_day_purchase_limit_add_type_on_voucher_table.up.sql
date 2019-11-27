ALTER TABLE vouchers
DROP COLUMN limit_per_user;

ALTER TABLE vouchers
DROP COLUMN day_purchase_limit;

ALTER TABLE vouchers
ADD COLUMN type SMALLINT;

