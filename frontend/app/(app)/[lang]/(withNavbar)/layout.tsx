import "@/app/globals.css";
import { Navbar } from "./_components/header/navbar/Navbar";
import { Footer } from "./_components/footer/Footer";
import { BackToAdminBanner } from "./_components/header/login/BackToAdmin/BackToAdminBanner";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../_components/providers/AppProviders";

export default async function WithNavbarLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const { lang } = await params;
  const messages = await getMessages({ locale: lang });
  return (
    <>
      <div className="print:hidden sticky top-0 z-40">
        <AppProviders lang={lang} messages={messages}>
          <BackToAdminBanner />
        </AppProviders>
        <Navbar lang={lang} />
      </div>
      <div className="flex-1">{children}</div>
      <div className="print:hidden">
        <Footer lang={lang} />
      </div>
    </>
  );
}
