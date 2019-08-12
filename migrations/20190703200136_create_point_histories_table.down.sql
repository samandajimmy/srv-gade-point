DROP INDEX index_point_histories;
DROP TABLE IF EXISTS point_histories:

ALTER TABLE campaign_transactions
    ADD CONSTRAINT campaign_transactions_voucher_code_id_fkey FOREIGN KEY (voucher_code_id) REFERENCES voucher_codes(id);