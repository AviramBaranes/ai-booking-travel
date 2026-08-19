import type { Metadata } from "next";
import { getLang } from "@/shared/lang/lang";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../../_components/providers/AppProviders";

/**
 * The funnel (results/plans/order) is query-param driven — indexing it would
 * surface stale prices and empty search states. Inherited by every child page,
 * none of which define their own metadata.
 */
export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default async function BookingLayout({
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
