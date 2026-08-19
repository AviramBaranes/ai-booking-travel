import type { Metadata } from "next";
import { getLang } from "@/shared/lang/lang";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../../_components/providers/AppProviders";

/** Tokenized price offers — private links, never indexable. */
export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default async function OfferLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const lang = await getLang();
  const messages = await getMessages({ locale: lang });

  return (
    <AppProviders lang={lang} messages={messages} showDevtools>
      {children}
    </AppProviders>
  );
}
