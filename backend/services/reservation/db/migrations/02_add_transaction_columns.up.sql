ALTER TABLE reservations
    ADD COLUMN payment_confirmation_code TEXT,
    ADD COLUMN payment_doc_num TEXT;