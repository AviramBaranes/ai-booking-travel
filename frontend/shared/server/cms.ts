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

export async function fetchBookingSettings(lang: string) {
  const payload = await getCachedPayload();
  return payload.findGlobal({
    slug: "booking-settings",
    draft: false,
    locale: lang === "en" ? "en" : "he",
  });
}

let payloadPromise: ReturnType<typeof getPayload> | null = null;
export function getCachedPayload() {
  if (!payloadPromise) {
    payloadPromise = getPayload({ config });
  }

  return payloadPromise;
}