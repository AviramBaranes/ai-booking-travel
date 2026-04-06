"use client";

import { useSearchParams, useRouter, usePathname } from "next/navigation";
import { useCallback } from "react";

export function useUrlFilters<K extends string>(
  keys: K[],
): [Record<K, string>, (updates: Partial<Record<K, string>>) => void] {
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();

  const filters = {} as Record<K, string>;
  for (const key of keys) {
    filters[key] = searchParams.get(key) ?? "";
  }

  const setFilters = useCallback(
    (updates: Partial<Record<K, string>>) => {
      const params = new URLSearchParams(searchParams.toString());
      for (const [key, value] of Object.entries(updates)) {
        if (value) {
          params.set(key, value as string);
        } else {
          params.delete(key);
        }
      }
      const qs = params.toString();
      router.replace(qs ? `${pathname}?${qs}` : pathname, { scroll: false });
    },
    [searchParams, router, pathname],
  );

  return [filters, setFilters];
}
