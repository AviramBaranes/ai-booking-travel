import Image from "next/image";
import Link from "next/link";

export function Logo({ lang, className }: { lang: string; className?: string }) {
  return (
    <Link href={`/${lang}`} className={className}>
      <Image
        src="/logo.png"
        alt="AIBookingTravel"
        width={168}
        height={32}
        className="object-contain"
        priority
      />
    </Link>
  );
}
