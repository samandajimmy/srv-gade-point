ALTER TABLE campaigns
ADD COLUMN metadata jsonb NOT NULL default '{}'::jsonb;