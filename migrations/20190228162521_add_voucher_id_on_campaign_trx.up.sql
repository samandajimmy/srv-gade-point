ALTER TABLE campaign_transactions
ADD COLUMN voucher_id SMALLINT REFERENCES vouchers(id) NULL,
ALTER COLUMN campaign_id DROP NOT NULL,
ALTER COLUMN campaign_id SET DEFAULT NULL;
