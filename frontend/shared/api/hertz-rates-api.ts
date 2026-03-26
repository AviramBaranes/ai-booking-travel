import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function listHertzMarkupRates(
  params: booking.ListHertzMarkupRatesRequest,
) {
  return withErrorHandler((client) =>
    client.booking.ListHertzMarkupRates(params),
  );
}

export function createHertzMarkupRate(
  data: booking.CreateHertzMarkupRateRequest,
) {
  return withErrorHandler((client) =>
    client.booking.CreateHertzMarkupRate(data),
  );
}

export function updateHertzMarkupRate(
  id: number,
  data: booking.UpdateHertzMarkupRateRequest,
) {
  return withErrorHandler((client) =>
    client.booking.UpdateHertzMarkupRate(id, data),
  );
}

export function deleteHertzMarkupRate(id: number) {
  return withErrorHandler((client) => client.booking.DeleteHertzMarkupRate(id));
}
