import { translation } from "../client";
import { withErrorHandler } from "./_api";

export function listBrokerTranslations(
  params: translation.ListBrokerTranslationsParams,
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
  data: translation.UpdateBrokerTranslationParams,
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
