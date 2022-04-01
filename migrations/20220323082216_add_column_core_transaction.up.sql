DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_payment') THEN

        CREATE TYPE status_payment AS ENUM (
            'success',
            'failed_inquiry',
            'failed_payment'
        );

    END IF;    
END
$$;

ALTER TABLE core_transactions
DROP COLUMN total_reward,
ALTER COLUMN transaction_type TYPE VARCHAR(4),
ADD COLUMN root_ref_trx VARCHAR(50),
ADD COLUMN inq_status status_payment default null;