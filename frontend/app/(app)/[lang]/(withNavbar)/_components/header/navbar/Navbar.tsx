import { AppProviders } from "../../../../_components/providers/AppProviders";
import { getMessages } from "next-intl/server";
import { NavbarActions } from "./NavbarActions";
import { NavbarLinks } from "./NavbarLinks";

interface NavbarProps {
  lang: string;
}

export async function Navbar({ lang }: NavbarProps) {
  const messages = await getMessages({ locale: lang });
  return (
    <header className="bg-white shadow-card">
      <nav className="mx-auto flex h-20 w-11/12 items-center justify-between px-6">
        <NavbarLinks lang={lang} />

        <AppProviders showDevtools={false} lang={lang} messages={messages}>
            <NavbarActions />
        </AppProviders>
      </nav>
    </header>
  );
}
