-- Table: point_histories
-- for status 0 --> pending
--            1 --> success

CREATE TABLE IF NOT EXISTS point_histories (
    id SERIAL PRIMARY KEY NOT NULL,
    cif VARCHAR(50) NOT NULL,
    point_amount INTEGER NOT NULL,
    transaction_type VARCHAR(2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    used_for VARCHAR NULL,
    ref_core VARCHAR NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    reward_id INTEGER REFERENCES rewards(id),
    voucher_code_id INTEGER REFERENCES voucher_codes(id),
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX index_point_histories ON point_histories (cif, used_for, transaction_type, transaction_date, ref_core);

ALTER TABLE campaign_transactions DROP CONSTRAINT campaign_transactions_voucher_code_id_fkey;
