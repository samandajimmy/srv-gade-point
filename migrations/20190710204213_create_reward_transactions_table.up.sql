-- Table: reward_transactions
-- for status        0 --> inquiry
--                   1 --> Success
--                   2 --> Reject       


CREATE TABLE IF NOT EXISTS reward_transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    ref_core VARCHAR(50),
    ref_id VARCHAR(50),
    reward_id VARCHAR(50),
    cif VARCHAR(50),
    used_promo_code VARCHAR(20),
    inquiry_date TIMESTAMP DEFAULT NULL,
    successed_date TIMESTAMP DEFAULT NULL,
    rejected_date TIMESTAMP DEFAULT NULL,
    timeout_date TIMESTAMP DEFAULT NULL,
    request_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL

);

CREATE INDEX index_reward_transactions ON reward_transactions (cif, used_promo_code, status);