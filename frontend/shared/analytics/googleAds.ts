import { IS_INDEXABLE } from "@/shared/seo/site";

export const GOOGLE_ADS_ID = "AW-18422971558";

/** The "רכישה" (purchase) conversion action in Google Ads. */
const PURCHASE_CONVERSION = `${GOOGLE_ADS_ID}/LbTaCOSzs_4cEKb54NBE`;

/**
 * Only production reports to Google Ads, so stage and local bookings never
 * show up as conversions.
 */
export const GOOGLE_ADS_ENABLED = IS_INDEXABLE;

declare global {
  interface Window {
    gtag?: (...args: unknown[]) => void;
  }
}

/**
 * Reports a paid booking. The reservation id doubles as the transaction id, so
 * Google Ads counts each reservation once even if this fires again.
 */
export function trackPurchaseConversion({
  reservationId,
  value,
  currency,
}: {
  reservationId: number;
  value: number;
  currency: string;
}) {
  if (!GOOGLE_ADS_ENABLED) return;

  window.gtag?.("event", "conversion", {
    send_to: PURCHASE_CONVERSION,
    value,
    currency,
    transaction_id: String(reservationId),
  });
}
