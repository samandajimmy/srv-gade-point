ALTER TABLE referral_transactions
ALTER COLUMN reward_referral TYPE INTEGER,
DROP CONSTRAINT referral_transactions_reward_id_fkey,
DROP COLUMN reward_id
;