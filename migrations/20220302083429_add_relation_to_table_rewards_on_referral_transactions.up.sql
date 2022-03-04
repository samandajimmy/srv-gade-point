ALTER TABLE referral_transactions
ALTER COLUMN reward_referral TYPE FLOAT,
ADD COLUMN reward_id INTEGER,
ADD CONSTRAINT referral_transactions_reward_id_fkey FOREIGN KEY (reward_id) REFERENCES rewards (id)
;