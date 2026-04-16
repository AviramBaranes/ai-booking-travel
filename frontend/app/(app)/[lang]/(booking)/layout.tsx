import { getLang } from "@/shared/lang/lang";
import { BookingProviders } from "./_components/BookingProviders";
import { getMessages } from "next-intl/server";

export default async function BookingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const lang = await getLang();
  const messages = await getMessages({ locale: lang });

  return (
    <BookingProviders lang={lang} messages={messages}>
      {children}
    </BookingProviders>
  );
}
