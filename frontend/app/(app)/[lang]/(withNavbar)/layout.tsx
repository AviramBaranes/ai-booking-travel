import "@/app/globals.css";
import { Navbar } from "./_components/header/navbar/Navbar";
import { Footer } from "./_components/footer/Footer";
import { BackToAdminBanner } from "./_components/header/login/BackToAdmin/BackToAdminBanner";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../_components/providers/AppProviders";
import { notFound } from "next/navigation";

export default async function WithNavbarLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const { lang } = await params;

  if (!["he", "en"].includes(lang)) {
    notFound();
  }

  const messages = await getMessages({ locale: lang });
  return (
    <div className="min-h-screen flex flex-col">
      <div className="print:hidden sticky top-0 z-40">
        <AppProviders lang={lang} messages={messages}>
          <BackToAdminBanner />
        </AppProviders>
        <Navbar lang={lang} />
      </div>

      <main className="flex-1 min-h-[70vh]">{children}</main>

      <div className="print:hidden">
        <Footer lang={lang} />
      </div>
    </div>
  );
}
