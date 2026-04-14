import { formatPrice } from "@/shared/utils/formatPrice";
import Image from "next/image";

interface PriceDetailRowProps {
  iconSrc: string;
  altText: string;
  label: string;
  price: number;
  currency: string;
}

export function PriceDetailRow({
  iconSrc,
  altText,
  label,
  price,
  currency,
}: PriceDetailRowProps) {
  return (
    <div className="flex justify-between items-center">
      <div className="flex gap-2 items-center">
        <Image
          src={iconSrc}
          alt={altText}
          width={24}
          height={24}
          className="w-6 h-6"
        />
        <p className="type-label text-navy font-normal">{label}</p>
      </div>
      <p className="type-label text-navy font-normal">
        {formatPrice(price, currency)}
      </p>
    </div>
  );
}
