import { APIError, booking_handlers, price_offer } from "../client";
import { withErrorHandler } from "./_api";

export function createPriceOffer(params: price_offer.CreatePriceOfferParams) {
  return withErrorHandler((client) => client.booking.CreatePriceOffer(params));
}

export function listPriceOffers(params: price_offer.ListPriceOffersRequest) {
  return withErrorHandler((client) => client.booking.ListPriceOffers(params));
}

export function getAgentPriceOffer(
  id: number,
  onNotFound?: (error: APIError) => void,
) {
  const options = onNotFound ? { onExpectedError: { 404: onNotFound } } : undefined;
  return withErrorHandler((client) => client.booking.GetAgentPriceOffer(id), options);
}

export function getClientPriceOffer(token: string, onNotFound?: (error: APIError) => void) {
  const options = onNotFound ? { onExpectedError: { 404: onNotFound } } : undefined;
  return withErrorHandler((client) =>
    client.booking.GetClientPriceOffer(token),
    options
  );
}

export function updatePriceOffer(
  id: number,
  params: price_offer.UpdatePriceOfferParams,
) {
  return withErrorHandler((client) =>
    client.booking.UpdatePriceOffer(id, params),
  );
}

export function renewPriceOffer(id: number) {
  return withErrorHandler((client) => client.booking.RenewPriceOffer(id));
}

export function bookPriceOffer(params: booking_handlers.BookPriceOfferParams) {
  return withErrorHandler((client) => client.booking.BookPriceOffer(params));
}

export function approvePriceOffer(id: number) {
  return withErrorHandler((client) => client.booking.ApprovePriceOffer(id));
}
