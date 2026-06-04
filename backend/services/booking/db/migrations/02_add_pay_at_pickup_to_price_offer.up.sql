ALTER TABLE price_offers
ADD COLUMN pay_at_pickup JSONB NOT NULL DEFAULT '{}'::jsonb;