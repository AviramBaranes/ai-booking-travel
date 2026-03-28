import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function listBrokerTranslations(
  params: booking.ListBrokerTranslationsRequest,
) {
  return withErrorHandler((client) =>
    client.booking.ListBrokerTranslations(params),
  );
}

export function deleteBrokerTranslation(id: number) {
  return withErrorHandler((client) =>
    client.booking.DeleteBrokerTranslation(id),
  );
}

export function updateBrokerTranslation(
  id: number,
  data: booking.UpdateBrokerTranslationRequest,
) {
  return withErrorHandler((client) =>
    client.booking.UpdateBrokerTranslation(id, data),
  );
}

export function verifyBrokerTranslation(id: number) {
  return withErrorHandler((client) =>
    client.booking.VerifyBrokerTranslation(id),
  );
}
