import { markup_rate } from "../client";
import { withErrorHandler } from "./_api";

export function listMarkupRates(
  params: markup_rate.ListMarkupRatesParams,
) {
  return withErrorHandler((client) =>
    client.booking.ListMarkupRates(params),
  );
}

export function createMarkupRate(
  data: markup_rate.CreateMarkupRateParams,
) {
  return withErrorHandler((client) =>
    client.booking.CreateMarkupRate(data),
  );
}

export function updateMarkupRate(
  id: number,
  data: markup_rate.UpdateMarkupRateParams,
) {
  return withErrorHandler((client) =>
    client.booking.UpdateMarkupRate(id, data),
  );
}

export function deleteMarkupRate(id: number) {
  return withErrorHandler((client) => client.booking.DeleteMarkupRate(id));
}
