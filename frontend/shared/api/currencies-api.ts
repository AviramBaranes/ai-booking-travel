import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function listCurrencies() {
  return withErrorHandler((client) => client.booking.ListCurrencies());
}

export function createCurrency(data: booking.CreateCurrencyRequest) {
  return withErrorHandler((client) => client.booking.CreateCurrency(data));
}

export function updateCurrency(
  id: number,
  data: booking.UpdateCurrencyRequest,
) {
  return withErrorHandler((client) => client.booking.UpdateCurrency(id, data));
}

export function deleteCurrency(id: number) {
  return withErrorHandler((client) => client.booking.DeleteCurrency(id));
}
