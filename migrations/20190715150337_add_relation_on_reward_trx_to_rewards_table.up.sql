ALTER TABLE reward_transactions
RENAME inquiry_date TO inquired_date;

ALTER TABLE reward_transactions
RENAME successed_date TO succeeded_date;

ALTER TABLE reward_transactions
ADD COLUMN transaction_date TIMESTAMP DEFAULT NULL,
ALTER COLUMN reward_id TYPE INTEGER USING (trim(reward_id)::INTEGER),
ADD CONSTRAINT reward_transactions_reward_id_fkey FOREIGN KEY (reward_id) REFERENCES rewards (id);