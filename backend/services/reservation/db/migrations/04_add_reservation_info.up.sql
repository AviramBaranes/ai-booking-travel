CREATE TABLE supplier_terms(
    id BIGSERIAL PRIMARY KEY,
    broker broker NOT NULL,
    supplier_code TEXT NOT NULL,
    pickup_location_id TEXT NOT NULL,
    terms JSON NOT NULL,

    UNIQUE(broker, supplier_code, pickup_location_id)
);

ALTER TABLE reservations
ADD COLUMN excess INT NOT NULL DEFAULT 0,
ADD COLUMN excess_currency TEXT NOT NULL DEFAULT '',
ADD COLUMN supplier_terms_id BIGINT REFERENCES supplier_terms(id),
ADD COLUMN pickup_details JSON NOT NULL DEFAULT '{}',
ADD COLUMN dropoff_details JSON,
ADD COLUMN pickup_location_code TEXT NOT NULL DEFAULT '',
ADD COLUMN dropoff_location_code TEXT NOT NULL DEFAULT '',
ADD COLUMN supplier_paid_at TIMESTAMPTZ,
ADD COLUMN supplier_expense_id TEXT,
ADD COLUMN payment_received_at TIMESTAMPTZ; -- For credit card payment, the date when the payment from the credit card settled at the bank.

CREATE INDEX idx_reservations_supplier_terms_id ON reservations (supplier_terms_id);
