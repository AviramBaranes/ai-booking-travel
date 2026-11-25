-- The driver's mobile phone, entered on the order page. NULL for reservations made before it was collected.
ALTER TABLE reservations ADD COLUMN driver_phone TEXT;
