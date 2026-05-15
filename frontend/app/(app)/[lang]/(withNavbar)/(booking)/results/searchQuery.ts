import { booking } from "@/shared/client";

export interface SearchQuery {
  pickupLocationId: number;
  dropoffLocationId: number;
  pickupDate: Date;
  dropoffDate: Date;
  pickupTime: string;
  dropoffTime: string;
  driverAge: number;
  couponCode?: string;
}

export function toSearchRequest(
  query: SearchQuery,
): booking.SearchAvailabilityRequest {
  const fmt = (d: Date) => d.toISOString().split("T")[0];
  return {
    PickupLocationID: query.pickupLocationId,
    DropoffLocationID: query.dropoffLocationId,
    PickupDate: fmt(query.pickupDate),
    DropoffDate: fmt(query.dropoffDate),
    PickupTime: query.pickupTime,
    DropoffTime: query.dropoffTime,
    DriverAge: query.driverAge,
    CouponCode: query.couponCode ?? "",
  };
}

export function searchRequestToParams(
  request: booking.SearchAvailabilityRequest,
): string {
  const params = new URLSearchParams({
    pl: String(request.PickupLocationID),
    dl: String(request.DropoffLocationID),
    pd: request.PickupDate,
    pt: request.PickupTime,
    dd: request.DropoffDate,
    dt: request.DropoffTime,
    da: String(request.DriverAge),
  });

  if (request.CouponCode) {
    params.set("cc", request.CouponCode);
  }

  return params.toString();
}

export function parseSearchQuery(params: URLSearchParams): SearchQuery | null {
  const pl = params.get("pl");
  const dl = params.get("dl");
  const pd = params.get("pd");
  const pt = params.get("pt");
  const dd = params.get("dd");
  const dt = params.get("dt");
  const da = params.get("da");

  if (!pl || !dl || !pd || !pt || !dd || !dt || !da) return null;

  const pickupLocationId = parseInt(pl, 10);
  const dropoffLocationId = parseInt(dl, 10);
  const driverAge = parseInt(da, 10);
  const pickupDate = new Date(pd);
  const dropoffDate = new Date(dd);

  if (
    isNaN(pickupLocationId) ||
    isNaN(dropoffLocationId) ||
    isNaN(driverAge) ||
    isNaN(pickupDate.getTime()) ||
    isNaN(dropoffDate.getTime())
  ) {
    return null;
  }

  return {
    pickupLocationId,
    dropoffLocationId,
    pickupDate,
    dropoffDate,
    pickupTime: pt,
    dropoffTime: dt,
    driverAge,
    couponCode: params.get("cc") ?? undefined,
  };
}
