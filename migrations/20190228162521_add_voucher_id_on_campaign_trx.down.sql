ALTER TABLE campaign_transactions
DROP COLUMN voucher_id,
ALTER COLUMN campaign_id DROP DEFAULT,
ALTER COLUMN campaign_id SET NOT NULL;