import { AddonsGallery } from "@/payload-types";
import { useSuspenseQuery } from "@tanstack/react-query";

export const addonsGalleryKey = ["cms", "addonsGallery"] as const;

export function useAddonsGallery() {
  return useSuspenseQuery<AddonsGallery>({
    queryKey: addonsGalleryKey,
    queryFn: async () => {
      const res = await fetch("/api/globals/addonsGallery");
      if (!res.ok) throw new Error("Failed to fetch AddonsGallery");
      return res.json();
    },
  });
}
