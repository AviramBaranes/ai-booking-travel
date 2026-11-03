-- The new enum value is not used in this transaction, so nothing else needs to be split out.
ALTER TYPE broker ADD VALUE IF NOT EXISTS 'hertz_gr';
