import { getLang } from "@/shared/lang/lang";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../_components/providers/AppProviders";

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
