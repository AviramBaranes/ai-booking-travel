-- name: CreatePriceOffer :one
INSERT INTO price_offers (
    agent_id,
    name,
    pickup_location_id,
    dropoff_location_id,
    pickup_date,
    return_date,
    rental_days,
    pickup_time,
    dropoff_time,
    driver_age,
    rate_qualifier,
    supplier_code,
    car_details,
    plan_inclusions,
    currency_code,
    purchase_price,
    markup_percentage,
    broker_erp_price,
    bt_erp_price,
    total_price,
    offered_currency_code,
    offered_price
) VALUES (
    sqlc.arg(agent_id),
    sqlc.arg(name),
    sqlc.arg(pickup_location_id),
    sqlc.arg(dropoff_location_id),
    sqlc.arg(pickup_date),
    sqlc.arg(return_date),
    sqlc.arg(rental_days),
    sqlc.arg(pickup_time),
    sqlc.arg(dropoff_time),
    sqlc.arg(driver_age),
    sqlc.arg(rate_qualifier),
    sqlc.arg(supplier_code),
    sqlc.arg(car_details),
    sqlc.arg(plan_inclusions),
    sqlc.arg(currency_code),
    sqlc.arg(purchase_price),
    sqlc.arg(markup_percentage),
    sqlc.arg(broker_erp_price),
    sqlc.arg(bt_erp_price),
    sqlc.arg(total_price),
    sqlc.arg(offered_currency_code),
    sqlc.arg(offered_price)
)
RETURNING *;   

-- name: GetPriceOfferByToken :one
SELECT price_offers.* , pl.name AS pickup_location, dl.name AS dropoff_location
FROM price_offers
    JOIN locations pl ON price_offers.pickup_location_id = pl.id
    JOIN locations dl ON price_offers.dropoff_location_id = dl.id
 WHERE token = sqlc.arg(token) AND status != 'unavailable';

-- name: GetPriceOfferById :one
SELECT price_offers.* , pl.name AS pickup_location, dl.name AS dropoff_location
FROM price_offers
    JOIN locations pl ON price_offers.pickup_location_id = pl.id
    JOIN locations dl ON price_offers.dropoff_location_id = dl.id
WHERE price_offers.id = sqlc.arg(id) AND price_offers.agent_id = sqlc.arg(agent_id);

-- name: UpdatePriceOffer :exec
UPDATE price_offers SET
    status = COALESCE(sqlc.narg(status), status),
    name = COALESCE(sqlc.narg(name), name),
    offered_currency_code = COALESCE(sqlc.narg(offered_currency_code), offered_currency_code),
    offered_price = COALESCE(sqlc.narg(offered_price), offered_price),
    updated_at = now()
WHERE id = sqlc.arg(id) AND agent_id = sqlc.arg(agent_id) AND status != 'unavailable';

-- name: RenewPriceOfferDetails :exec
UPDATE price_offers SET
    car_details = sqlc.arg(car_details),
    plan_inclusions = sqlc.arg(plan_inclusions),
    currency_code = sqlc.arg(currency_code),
    purchase_price = sqlc.arg(purchase_price),
    markup_percentage = sqlc.arg(markup_percentage),
    broker_erp_price = sqlc.arg(broker_erp_price),
    bt_erp_price = sqlc.arg(bt_erp_price),
    total_price = sqlc.arg(total_price),
    renewed_at = now(),
    updated_at = now()
WHERE id = sqlc.arg(id) AND agent_id = sqlc.arg(agent_id) AND status != 'unavailable';

-- name: RenewPriceOfferUnavailable :exec
UPDATE price_offers SET
    status = 'unavailable',
    renewed_at = now(),
    updated_at = now()
WHERE id = sqlc.arg(id) AND agent_id = sqlc.arg(agent_id);

-- name: SetPriceOfferRenewedAt :exec
UPDATE price_offers SET
    renewed_at = sqlc.arg(renewed_at)
WHERE id = sqlc.arg(id) AND agent_id = sqlc.arg(agent_id);

-- name: ListPriceOffersByAgent :many
SELECT price_offers.id, status, price_offers.name, 
pl.name AS pickup_location, dl.name AS dropoff_location, 
pickup_date, return_date, pickup_time, dropoff_time, 
currency_code, total_price, offered_currency_code, offered_price, 
price_offers.created_at
FROM price_offers
    JOIN locations pl ON price_offers.pickup_location_id = pl.id
    JOIN locations dl ON price_offers.dropoff_location_id = dl.id
WHERE agent_id = sqlc.arg(agent_id)
  AND (sqlc.narg(status)::offer_status IS NULL OR status = sqlc.narg(status)::offer_status)
  AND (sqlc.narg(name_search)::text IS NULL OR price_offers.name ILIKE '%' || sqlc.narg(name_search)::text || '%')
ORDER BY price_offers.created_at DESC
LIMIT sqlc.arg(page_size)
OFFSET sqlc.arg(page_offset);

-- name: CountPriceOffersByAgent :one
SELECT COUNT(*)::BIGINT AS total
FROM price_offers
WHERE agent_id = sqlc.arg(agent_id)
  AND (sqlc.narg(status)::offer_status IS NULL OR status = sqlc.narg(status)::offer_status)
  AND (sqlc.narg(name_search)::text IS NULL OR price_offers.name ILIKE '%' || sqlc.narg(name_search)::text || '%');