CREATE TYPE penalty_type AS ENUM ('cancellation', 'no_show');

CREATE TABLE
    reservation_penalties (
        id BIGSERIAL PRIMARY KEY,
        reservation_id BIGINT NOT NULL REFERENCES reservations (id),
        penalty_type penalty_type NOT NULL,
        created_by_user_id BIGINT,

        currency_code TEXT NOT NULL,
        currency_rate NUMERIC(12, 4) NOT NULL,
        amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0), -- The supplier charges us and we charge the customer the same amount.

        paid_at TIMESTAMPTZ,
        invoice_doc_num TEXT,
        payment_doc_num TEXT,

        supplier_paid_at TIMESTAMPTZ,
        supplier_expense_id TEXT,

        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),

        -- A cancellation and a no-show are mutually exclusive, so a reservation carries at most one penalty.
        UNIQUE (reservation_id)
    );

CREATE INDEX idx_reservation_penalties_unpaid ON reservation_penalties (created_at) WHERE paid_at IS NULL;

CREATE INDEX idx_reservation_penalties_unpaid_supplier ON reservation_penalties (created_at) WHERE supplier_paid_at IS NULL;
