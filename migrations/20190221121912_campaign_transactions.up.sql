    -- Table: campaign_transactions
    -- for transaction_type D --> debet and K --> kredit
    CREATE TABLE IF NOT EXISTS campaign_transactions (
        id SERIAL PRIMARY KEY NOT NULL,
        user_id VARCHAR(50) NOT NULL,
        point_amount SMALLINT NOT NULL,
        transaction_type CHAR(2) NOT NULL,
        transaction_date TIMESTAMP NOT NULL,
        campaign_id SMALLINT,
        promo_code_id SMALLINT,
        updated_at TIMESTAMP DEFAULT NULL,
        created_at TIMESTAMP DEFAULT NULL
    );

    CREATE INDEX index_campaign_transactions ON campaign_transactions (user_id, point_amount, transaction_type, transaction_date, campaign_id, promo_code_id);
