import "server-only";
import { getPayload } from "payload";
import config from "@payload-config";

export async function fetchSuppliersGallery() {
  const payload = await getPayload({ config });
  return payload.findGlobal({ slug: "suppliersGallery", draft: false });
}

export async function fetchAddonsGallery() {
  const payload = await getPayload({ config });
  return payload.findGlobal({ slug: "addonsGallery", draft: false });
}

export async function fetchBookingSettings() {
  const payload = await getPayload({ config });
  return payload.findGlobal({ slug: "booking-settings", draft: false });
}
