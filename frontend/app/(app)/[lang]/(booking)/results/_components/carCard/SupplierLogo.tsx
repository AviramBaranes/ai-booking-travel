import { SuppliersGallery } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Image from "next/image";

type Supplier = NonNullable<SuppliersGallery["suppliers"]>[number];
type SupplierMedia = Supplier["media"];

export function SupplierLogo({
  supplierName,
  supplierGallery,
}: {
  supplierName: string;
  supplierGallery: SuppliersGallery;
}) {
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
      width={80}
      height={40}
      className="w-25 h-auto"
    />
  );
}

function isPopulatedMedia(
  media: SupplierMedia,
): media is Populated<SupplierMedia> {
  return typeof media !== "number";
}
