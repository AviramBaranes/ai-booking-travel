import "server-only";
import { getPayload } from "payload";
import config from "@payload-config";

export async function fetchSuppliersGallery() {
  const payload = await getCachedPayload();
  return payload.findGlobal({ slug: "suppliersGallery", draft: false });
}

export async function fetchAddonsGallery() {
  const payload = await getCachedPayload();
  return payload.findGlobal({ slug: "addonsGallery", draft: false });
}

export async function fetchBookingSettings() {
  const payload = await getCachedPayload();
  return payload.findGlobal({ slug: "booking-settings", draft: false });
}

let payloadPromise: ReturnType<typeof getPayload> | null = null;
let initCount = 0;

export function getCachedPayload() {
  if (!payloadPromise) {
    initCount += 1;
    console.log(`[Payload] init #${initCount}`);

    payloadPromise = getPayload({ config });
  }

  return payloadPromise;
}