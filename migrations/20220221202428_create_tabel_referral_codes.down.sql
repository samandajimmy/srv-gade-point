DROP INDEX index_referral_codes;
DROP INDEX unique_index_referral_codes;

ALTER TABLE referral_codes
DROP CONSTRAINT fk_referral_codes;

DROP TABLE IF EXISTS referral_codes;