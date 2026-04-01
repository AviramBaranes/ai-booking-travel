import { Footer } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import Link from "next/link";

interface FooterThirdFloorProps {
  links: Footer["thirdFloorLinks"];
  rightsText: string;
  lang: string;
}

export function FooterThirdFloor({
  links,
  rightsText,
  lang,
}: FooterThirdFloorProps) {
  return (
    <div className="flex items-center justify-between bg-navy py-5 px-20">
      <div className="flex items-start gap-6">
        {links?.map((link) => (
          <Link
            key={link.id}
            className="text-white text-sm"
            href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
          >
            {link.label}
          </Link>
        ))}
      </div>
      <span className="text-white opacity-65 text-xs">{rightsText}</span>
    </div>
  );
}
