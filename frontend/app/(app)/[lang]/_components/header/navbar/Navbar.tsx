import Image from "next/image";
import { getPayload } from "payload";
import config from "@payload-config";
import Link from "next/link";
import { MegaDropdown } from "./MegaDropdown";
import type { Populated } from "@/shared/types/payload";
import { LogoutButton } from "./LogoutButton";
import { LangSwitcher } from "../login/LangSwitcher";
import { LoginModal } from "../login/LoginModal";
import { AppProviders } from "../../providers/AppProviders";
import { getMessages } from "next-intl/server";

async function getHeaderData(lang: string) {
  const payload = await getPayload({ config });
  return payload.findGlobal({
    slug: "header",
    locale: lang as "he" | "en",
    draft: false,
  });
}

interface NavbarProps {
  lang: string;
  isAuthenticated: boolean;
  // When true, hides the login/logout buttons and language switcher. Used 404 page outside of the [lang] path
  isRootLayout?: boolean;
}

export async function Navbar({
  lang,
  isAuthenticated,
  isRootLayout = false,
}: NavbarProps) {
  const headerData = await getHeaderData(lang);
  const messages = await getMessages({ locale: lang });
  return (
    <header className="sticky top-0 z-40 bg-white shadow-card">
      <nav className="mx-auto flex h-20 w-11/12 items-center justify-between px-6">
        <div className="flex items-center gap-8">
          <Link href={`/${lang}`}>
            <Image
              src="/logo.png"
              alt="AIBookingTravel"
              width={160}
              height={40}
              className="w-40 h-10"
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

        {!isRootLayout && (
          <AppProviders lang={lang} messages={messages}>
            <div className="flex items-center gap-4">
              <LangSwitcher lang={lang} />
              {isAuthenticated ? <LogoutButton /> : <LoginModal />}
            </div>
          </AppProviders>
        )}
      </nav>
    </header>
  );
}
