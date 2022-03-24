ALTER TABLE core_transactions
DROP COLUMN transaction_type;
DROP COLUMN total_reward;

CREATE TYPE status_enum AS ENUM('success', 'failed_inquiry', 'failed_payment');

ALTER TABLE core_transactions
ADD COLUMN inq_status TYPE status_enum using inq_status::status_enum,
ADD COLUMN root_ref_trx varchar(50),
ADD COLUMN transaction_type VARCHAR(2);