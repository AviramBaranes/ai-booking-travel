-- plan_inclusions holds what the user saw, translated when they searched in Hebrew. Vouchers are
-- always in English, so the English source is kept, only when it differs, and is NULL otherwise.
ALTER TABLE reservations ADD COLUMN plan_inclusions_en TEXT[];
