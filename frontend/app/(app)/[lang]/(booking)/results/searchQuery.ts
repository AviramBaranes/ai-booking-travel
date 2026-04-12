import { booking } from "@/shared/client";

export interface SearchQuery {
  pickupLocationId: number;
  returnLocationId: number;
  pickupDate: Date;
  returnDate: Date;
  pickupTime: string;
  returnTime: string;
  driverAge: number;
  couponCode?: string;
}

export function toSearchRequest(
  query: SearchQuery,
): booking.SearchAvailabilityRequest {
  const fmt = (d: Date) => d.toISOString().split("T")[0];
  return {
    PickupLocationID: query.pickupLocationId,
    DropoffLocationID: query.returnLocationId,
    PickupDate: fmt(query.pickupDate),
    DropoffDate: fmt(query.returnDate),
    PickupTime: query.pickupTime,
    DropoffTime: query.returnTime,
    DriverAge: query.driverAge,
    CouponCode: query.couponCode ?? "",
  };
}

export function searchRequestToParams(
  request: booking.SearchAvailabilityRequest,
): string {
  const params = new URLSearchParams({
    pl: String(request.PickupLocationID),
    rl: String(request.DropoffLocationID),
    pd: request.PickupDate,
    pt: request.PickupTime,
    rd: request.DropoffDate,
    rt: request.DropoffTime,
    da: String(request.DriverAge),
  });

  if (request.CouponCode) {
    params.set("cc", request.CouponCode);
  }

  return params.toString();
}

export function parseSearchQuery(params: URLSearchParams): SearchQuery | null {
  const pl = params.get("pl");
  const rl = params.get("rl");
  const pd = params.get("pd");
  const pt = params.get("pt");
  const rd = params.get("rd");
  const rt = params.get("rt");
  const da = params.get("da");

  if (!pl || !rl || !pd || !pt || !rd || !rt || !da) return null;

  const pickupLocationId = parseInt(pl, 10);
  const returnLocationId = parseInt(rl, 10);
  const driverAge = parseInt(da, 10);
  const pickupDate = new Date(pd);
  const returnDate = new Date(rd);

  if (
    isNaN(pickupLocationId) ||
    isNaN(returnLocationId) ||
    isNaN(driverAge) ||
    isNaN(pickupDate.getTime()) ||
    isNaN(returnDate.getTime())
  ) {
    return null;
  }

  return {
    pickupLocationId,
    returnLocationId,
    pickupDate,
    returnDate,
    pickupTime: pt,
    returnTime: rt,
    driverAge,
    couponCode: params.get("cc") ?? undefined,
  };
}
