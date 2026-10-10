import { location } from "../client";
import { withErrorHandler } from "./_api";

export function searchLocations(
  query: string,
): Promise<location.SearchLocationResponse | null | undefined> {
  return withErrorHandler((client) =>
    client.booking.SearchLocations({ Search: query }),
  );
}

export function listLocations(params: location.ListLocationsParams) {
  return withErrorHandler((client) => client.booking.ListLocations(params));
}

export function insertLocation(data: location.InsertLocationParams) {
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

export function bulkToggleLocations(ids: number[], enabled: boolean) {
  return withErrorHandler((client) =>
    client.booking.BulkToggleLocations({ ids, enabled }),
  );
}

export function toggleIsAirport(id: number, is_airport: boolean) {
  return withErrorHandler((client) =>
    client.booking.ToggleLocationIsAirport(id, { is_airport }),
  );
}

export function insertLocationAlias(data: location.InsertLocationAliasesParams) {
  return withErrorHandler((client) =>
    client.booking.InsertLocationAlias(data),
  );
}

export function listLocationMergeSuggestions() {
  return withErrorHandler((client) =>
    client.booking.ListLocationMergeSuggestions(),
  );
}

export function searchLocationsForMerge(search: string) {
  return withErrorHandler((client) =>
    client.booking.SearchLocationsForMerge({ Search: search }),
  );
}

export function mergeLocations(data: location.MergeLocationsParams) {
  return withErrorHandler((client) => client.booking.MergeLocations(data));
}
