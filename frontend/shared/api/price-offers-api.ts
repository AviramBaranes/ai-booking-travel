import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function createPriceOffer(params: booking.CreatePriceOfferParams) {
  return withErrorHandler((client) => client.booking.CreatePriceOffer(params));
}

export function listPriceOffers(params: booking.ListPriceOffersRequest) {
  return withErrorHandler((client) => client.booking.ListPriceOffers(params));
}

export function getAgentPriceOffer(id: number) {
  return withErrorHandler((client) => client.booking.GetAgentPriceOffer(id));
}

export function getClientPriceOffer(token: string) {
  return withErrorHandler((client) =>
    client.booking.GetClientPriceOffer(token),
  );
}

export function updatePriceOffer(
  id: number,
  params: booking.UpdatePriceOfferParams,
) {
  return withErrorHandler((client) =>
    client.booking.UpdatePriceOffer(id, params),
  );
}

export function renewPriceOffer(id: number) {
  return withErrorHandler((client) => client.booking.RenewPriceOffer(id));
}
