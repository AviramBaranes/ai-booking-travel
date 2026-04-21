import { getMessages } from "next-intl/server";
import { AppProviders } from "../_components/providers/AppProviders";
import { getLang } from "@/shared/lang/lang";

export default async function AccountLayout({
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
