-- The full SupplierInfo payload returned by the broker availability call is stored once per
-- snapshot (not per plan) so that terms and station details do not bloat every row in `plans`.
ALTER TABLE available_plans_snapshots
    ADD COLUMN suppliers_info JSON NOT NULL DEFAULT '[]';

-- Price offers are booked long after their snapshot is gone, so the supplier terms and station
-- details of the chosen plan are denormalized onto the offer itself.
ALTER TABLE price_offers
    ADD COLUMN supplier_terms JSON,
    ADD COLUMN pickup_details JSON,
    ADD COLUMN dropoff_details JSON;
