import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function searchLocations(
  query: string,
): Promise<booking.SearchLocationResponse | null | undefined> {
  return withErrorHandler((client) =>
    client.booking.SearchLocations({ Search: query }),
  );
}
