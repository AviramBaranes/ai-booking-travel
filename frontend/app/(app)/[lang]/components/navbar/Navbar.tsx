import Image from "next/image";
import { LangSwitcher } from "../LangSwitcher";
import { LoginModal } from "../LoginModal";
import { getPayload } from "payload";
import config from "@payload-config";
import Link from "next/link";
import { MegaDropdown } from "./MegaDropdown";
import type { Populated } from "@/shared/types/payload";

async function getHeaderData(lang: string) {
  const payload = await getPayload({ config });
  return payload.findGlobal({
    slug: "header",
    locale: lang as "he" | "en",
    draft: false,
  });
}

export async function Navbar({
  lang,
  isAuthenticated,
}: {
  lang: string;
  isAuthenticated: boolean;
}) {
  const headerData = await getHeaderData(lang);
  return (
    <header className="sticky top-0 z-40 bg-white">
      <nav className="mx-auto flex h-20 w-11/12 items-center justify-between px-6">
        <div className="flex items-center gap-8">
          <Link href={`/${lang}`}>
            <Image
              src="/logo.png"
              alt="AIBookingTravel"
              width={160}
              height={40}
            />
          </Link>

          {headerData.links?.map((link) =>
            link.type === "link" ? (
              <Link
                key={link.id}
                href={`/${lang}/${(link.page as Populated<typeof link.page>)?.slug ?? ""}`}
                className="text-lg font-bold text-navy"
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

        <div className="flex items-center gap-4">
          <LangSwitcher lang={lang} />
          {!isAuthenticated && <LoginModal />}
        </div>
      </nav>
    </header>
  );
}
