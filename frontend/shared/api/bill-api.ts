import { billing, isAPIError } from "../client";
import { withErrorHandler } from "./_api";

export class BillingFailedError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "BillingFailedError";
  }
}

export function isBillingFailedError(err: unknown): err is BillingFailedError {
  return err instanceof BillingFailedError;
}

export function bill(params: billing.BillParams) {
  return withErrorHandler(async (client) => {
    try {
      return await client.billing.Bill(params);
    } catch (err) {
      if (isAPIError(err) && err.code === "unknown" && err.message) {
        throw new BillingFailedError(err.message);
      }
      throw err;
    }
  });
}

export function getOrderPaymentIframe(reservationId: number, isIls: boolean) {
  return withErrorHandler((client) =>
    client.billing.GenerateOrderIframe({
      orderId: reservationId,
      isIls,
    }),
  );
}
