import { booking } from "../client";
import { withErrorHandler } from "./_api";

export function searchLocations(
  query: string,
): Promise<booking.SearchLocationResponse | null | undefined> {
  return withErrorHandler((client) =>
    client.booking.SearchLocations({ Search: query }),
  );
}

export function listLocations(params: booking.ListLocationsRequest) {
  return withErrorHandler((client) => client.booking.ListLocations(params));
}

export function insertLocation(data: booking.InsertLocationParams) {
  return withErrorHandler((client) => client.booking.InsertLocation(data));
}

export function deleteLocation(id: number) {
  return withErrorHandler((client) => client.booking.DeleteLocation(id));
}

export function toggleLocation(id: number, enabled: boolean) {
  return withErrorHandler((client) =>
    client.booking.ToggleLocation(id, { enabled }),
  );
}
