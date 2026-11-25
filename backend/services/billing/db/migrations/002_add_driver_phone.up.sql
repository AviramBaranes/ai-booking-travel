-- The driver's mobile phone, kept on the pending payment until the booking is made after payment.
ALTER TABLE pending_customer_payments ADD COLUMN driver_phone TEXT;
