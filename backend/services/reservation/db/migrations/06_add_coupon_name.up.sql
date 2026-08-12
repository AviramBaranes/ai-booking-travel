ALTER TABLE reservations
ADD COLUMN coupon_name TEXT NOT NULL DEFAULT ''; -- The name of the coupon the customer booked with, empty when no coupon was applied.
