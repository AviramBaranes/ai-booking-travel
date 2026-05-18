import { markup_rate } from "../client";
import { withErrorHandler } from "./_api";

export function listHertzMarkupRates(
  params: markup_rate.ListHertzMarkupRatesParams,
) {
  return withErrorHandler((client) =>
    client.booking.ListHertzMarkupRates(params),
  );
}

export function createHertzMarkupRate(
  data: markup_rate.CreateHertzMarkupRateParams,
) {
  return withErrorHandler((client) =>
    client.booking.CreateHertzMarkupRate(data),
  );
}

export function updateHertzMarkupRate(
  id: number,
  data: markup_rate.UpdateHertzMarkupRateParams,
) {
  return withErrorHandler((client) =>
    client.booking.UpdateHertzMarkupRate(id, data),
  );
}

export function deleteHertzMarkupRate(id: number) {
  return withErrorHandler((client) => client.booking.DeleteHertzMarkupRate(id));
}
