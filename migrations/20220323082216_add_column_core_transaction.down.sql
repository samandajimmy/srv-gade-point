ALTER TABLE core_transactions
ADD COLUMN total_reward DECIMAL,
ALTER COLUMN transaction_type TYPE SMALLINT NOT NULL DEFAULT 0,
DROP COLUMN root_ref_trx,
DROP COLUMN inq_status;

DROP TYPE status_payment;