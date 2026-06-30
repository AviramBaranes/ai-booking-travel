import { AppProviders } from "../../../../_components/providers/AppProviders";
import { getMessages } from "next-intl/server";
import { NavbarActions } from "./NavbarActions";
import { NavbarLinks } from "./NavbarLinks";
import { MobileMenuDrawer } from "./MobileMenuDrawer";
import { getCachedPayload } from "@/shared/server/cms";
import { Logo } from "./Logo";

interface NavbarProps {
  lang: string;
}

export async function getHeaderData(lang: string) {
  const payload = await getCachedPayload();
  return payload.findGlobal({
    slug: "header",
    locale: lang as "he" | "en",
    draft: false,
  });
}

export async function Navbar({ lang }: NavbarProps) {
  const messages = await getMessages({ locale: lang });
  const headerData = await getHeaderData(lang);
  return (
    <header className="bg-white shadow-card">
      <nav className="mx-auto flex h-15 lg:h-20 w-11/12 items-center justify-between px-2 lg:px-6">
        <Logo lang={lang} className="lg:hidden w-1/3" />
        <NavbarLinks
          lang={lang}
          headerData={headerData}
          className="hidden lg:flex"
        />

        <div className="flex gap-3 items-center">
          <AppProviders showDevtools={false} lang={lang} messages={messages}>
            <NavbarActions />
          </AppProviders>
          <MobileMenuDrawer lang={lang} links={headerData.links ?? []} />
        </div>
      </nav>
    </header>
  );
}
