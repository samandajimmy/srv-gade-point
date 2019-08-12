ALTER TABLE campaigns
    ADD COLUMN type SMALLINT DEFAULT 0,
    ADD COLUMN validators JSONB NOT NULL;

ALTER TABLE campaign_transactions RENAME COLUMN ref_core TO reff_core;

ALTER TABLE point_histories
    DROP COLUMN ref_id;

ALTER TABLE campaign_transactions
    DROP COLUMN ref_id;