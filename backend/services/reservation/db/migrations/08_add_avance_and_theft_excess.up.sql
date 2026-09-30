-- The new enum value is not used in this transaction, so the columns can be added alongside it.
ALTER TYPE broker ADD VALUE IF NOT EXISTS 'avance';

ALTER TABLE reservations
ADD COLUMN theft_excess INT NOT NULL DEFAULT 0,
ADD COLUMN theft_excess_currency TEXT NOT NULL DEFAULT '';
