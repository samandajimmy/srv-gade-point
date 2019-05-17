ALTER TABLE campaign_transactions
ALTER COLUMN promo_code_id TYPE SMALLINT,
ALTER COLUMN campaign_id TYPE SMALLINT,
ALTER COLUMN transaction_type TYPE CHAR(2);