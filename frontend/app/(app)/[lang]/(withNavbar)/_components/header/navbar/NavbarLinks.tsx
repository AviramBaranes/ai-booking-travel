import Image from "next/image";
import Link from "next/link";
import { MegaDropdown } from "./MegaDropdown";
import type { Populated } from "@/shared/types/payload";
import { Header } from "@/payload-types";
import { Logo } from "./Logo";

interface NavbarLinksProps {
  lang: string;
  headerData: Header
  className?: string;
}

export async function NavbarLinks({ lang, headerData, className }: NavbarLinksProps) {
  return (
    <div className={`flex items-center gap-8 ${className}`}>
      <Logo lang={lang} />

      {headerData.links?.map((link) =>
        link.type === "link" ? (
          <Link
            key={link.id}
            href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
            className="type-h6 text-navy"
          >
            {link.label}
          </Link>
        ) : (
          <MegaDropdown
            key={link.id}
            label={link.megaLabel!}
            links={link.megaLinks ?? []}
            lang={lang}
          />
        ),
      )}
    </div>
  );
}
