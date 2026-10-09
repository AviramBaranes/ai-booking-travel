import { useMemo, useSyncExternalStore } from "react";

const STORAGE_KEY = "lastSearch";

interface SearchedLocation {
  id: number;
  name: string;
}

/** The last submitted search, kept in this browser so the form can be refilled. */
export interface LastSearch {
  pickupLocation: SearchedLocation;
  dropoffLocation: SearchedLocation;
  pickupDate: string; // yyyy-MM-dd
  pickupTime: string;
  dropoffDate: string; // yyyy-MM-dd
  dropoffTime: string;
  driverAge: number;
  couponCode?: string;
}

export function saveLastSearch(search: LastSearch) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(search));
  } catch {
    // Storage can be blocked or full; refilling the form is a convenience.
  }
}

function readRawLastSearch(): string | null {
  try {
    return localStorage.getItem(STORAGE_KEY);
  } catch {
    return null;
  }
}

function parseLastSearch(raw: string | null): LastSearch | null {
  if (!raw) return null;
  try {
    return JSON.parse(raw) as LastSearch;
  } catch {
    return null;
  }
}

// Saves happen right before navigating away, so there is nothing to subscribe to.
const subscribe = () => () => {};

/**
 * Reads the last search. It is null while rendering on the server, so `raw`
 * changes once the client hydrates — key a form on it to pick up the values.
 */
export function useLastSearch() {
  const raw = useSyncExternalStore(subscribe, readRawLastSearch, () => null);
  const lastSearch = useMemo(() => parseLastSearch(raw), [raw]);

  return { lastSearch, raw };
}
