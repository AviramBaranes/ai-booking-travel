import Image from "next/image";
import Link from "next/link";
import { MegaDropdown } from "./MegaDropdown";
import type { Populated } from "@/shared/types/payload";
import { getCachedPayload } from "@/shared/server/cms";

async function getHeaderData(lang: string) {
  const payload = await getCachedPayload();
  return payload.findGlobal({
    slug: "header",
    locale: lang as "he" | "en",
    draft: false,
  });
}

interface NavbarLinksProps {
  lang: string;
}

export async function NavbarLinks({ lang }: NavbarLinksProps) {
  const headerData = await getHeaderData(lang);
  return (
    <div className="flex items-center gap-8">
      <Link href={`/${lang}`}>
        <Image
          src="/logo.png"
          alt="AIBookingTravel"
          width={168}
          height={32}
          className="object-contain"
          priority
        />
      </Link>

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
