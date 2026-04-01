import { Footer } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Link from "next/link";

interface FooterFirstFloorProps {
  links: Footer["firstFloorLinks"];
  lang: string;
}

export function FooterFirstFloor({ links, lang }: FooterFirstFloorProps) {
  return (
    <div className="flex items-start gap-6 flex-wrap bg-dark-navy py-6 px-20">
      {links?.map((link) => (
        <Link
          key={link.id}
          className="text-white opacity-65 text-sm"
          href={`/${lang}/${(link.page as Populated<typeof link.page>).slug}`}
        >
          {link.label}
        </Link>
      ))}
    </div>
  );
}
