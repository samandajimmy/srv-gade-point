    ALTER TABLE referral_transactions
    DROP COLUMN trx_amount,
    DROP COLUMN loan_amount,
    DROP COLUMN interest_amount,
    DROP COLUMN trx_date,
    DROP COLUMN product_code,
    DROP COLUMN trx_type;

    -- Table: core_transactions
    -- for transaction_type 0 --> scheduler and 1 --> non scheduler
    CREATE TABLE IF NOT EXISTS core_transactions (
        id SERIAL PRIMARY KEY NOT NULL,
        transaction_id VARCHAR(50),
		transaction_amount DECIMAL,
		loan_amount DECIMAL,
		interest_amount DECIMAL,
		transaction_date TIMESTAMP,
		product_code VARCHAR(2),
		marketing_code VARCHAR(50),
		total_reward DECIMAL,
		transaction_type SMALLINT NOT NULL DEFAULT 0,
        updated_at TIMESTAMP DEFAULT NULL,
        created_at TIMESTAMP DEFAULT NULL
    );