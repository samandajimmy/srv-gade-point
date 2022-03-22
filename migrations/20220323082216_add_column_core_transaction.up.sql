ALTER TABLE core_transactions
DROP COLUMN trx_type;

ALTER TABLE core_transactions
ADD COLUMN inq_status SMALLINT NOT NULL DEFAULT 0,
ADD COLUMN root_ref_trx varchar(50),
ADD COLUMN trx_type VARCHAR(2);