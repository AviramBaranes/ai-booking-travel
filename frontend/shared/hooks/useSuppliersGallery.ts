import { SuppliersGallery } from "@/payload-types";
import { useSuspenseQuery } from "@tanstack/react-query";

export const suppliersGalleryKey = ["cms", "suppliersGallery"] as const;

export function useSupplierGallery() {
  return useSuspenseQuery<SuppliersGallery>({
    queryKey: suppliersGalleryKey,
    queryFn: async () => {
      const res = await fetch("/api/globals/suppliersGallery");
      if (!res.ok) throw new Error("Failed to fetch SuppliersGallery");
      return res.json();
    },
  });
}
