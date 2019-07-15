ALTER TABLE reward_transactions
RENAME inquired_date TO inquiry_date;

ALTER TABLE reward_transactions
RENAME succeeded_date TO successed_date;

ALTER TABLE reward_transactions
DROP COLUMN transaction_date,
ALTER COLUMN reward_id TYPE VARCHAR(50),
DROP CONSTRAINT reward_transactions_reward_id_fkey;