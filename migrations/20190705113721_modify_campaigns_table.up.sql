ALTER TABLE campaigns
    DROP COLUMN type,
    DROP COLUMN validators;

ALTER TABLE campaign_transactions RENAME reff_core TO ref_core;

ALTER TABLE point_histories
    ADD COLUMN ref_id VARCHAR;

ALTER TABLE campaign_transactions
    ADD COLUMN ref_id VARCHAR;