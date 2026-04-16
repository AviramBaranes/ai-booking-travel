import { Populated } from "@/shared/types/payload";
import { useSupplierGallery } from "@/shared/hooks/useSuppliersGallery";
import { SuppliersGallery } from "@/payload-types";
import Image from "next/image";

type Supplier = NonNullable<SuppliersGallery["suppliers"]>[number];
type SupplierMedia = Supplier["media"];

export function SupplierLogo({ supplierName }: { supplierName: string }) {
  const { data: supplierGallery } = useSupplierGallery();
  const supplier = supplierGallery.suppliers?.find(
    (s) => s.name.trim() === supplierName.trim(),
  );

  if (!supplier) return <h6>{supplierName}</h6>;
  const media = supplier.media;

  if (!isPopulatedMedia(media) || !media.url) {
    return <h6>{supplierName}</h6>;
  }

  return (
    <Image
      src={media.url}
      alt={supplier.name}
      width={112}
      height={40}
      className="w-28 h-10 object-cover"
    />
  );
}

function isPopulatedMedia(
  media: SupplierMedia,
): media is Populated<SupplierMedia> {
  return typeof media !== "number";
}
