ALTER TABLE vouchers ALTER COLUMN end_date SET NOT NULL;

ALTER TABLE vouchers
DROP COLUMN generator_type CASCADE;