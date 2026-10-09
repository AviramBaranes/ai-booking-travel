-- plan_inclusions holds what the agent saw, translated when they searched in Hebrew. The English
-- source is kept for the voucher, only when it differs, and is NULL otherwise.
ALTER TABLE price_offers ADD COLUMN plan_inclusions_en TEXT[];
