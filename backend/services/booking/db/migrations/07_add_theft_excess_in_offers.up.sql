ALTER TABLE price_offers
    ADD COLUMN theft_excess INT NOT NULL DEFAULT 0,
    ADD COLUMN theft_excess_currency TEXT NOT NULL DEFAULT '';
