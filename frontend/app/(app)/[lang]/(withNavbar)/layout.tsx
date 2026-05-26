import "@/app/globals.css";
import { Navbar } from "./_components/header/navbar/Navbar";
import { Footer } from "./_components/footer/Footer";
import { BackToAdminBanner } from "./_components/header/login/BackToAdmin/BackToAdminBanner";

export default async function WithNavbarLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const { lang } = await params;
  return (
    <>
      <div className="print:hidden sticky top-0 z-40">
        <BackToAdminBanner />
        <Navbar lang={lang} />
      </div>
      <div className="flex-1">{children}</div>
      <div className="print:hidden">
        <Footer lang={lang} />
      </div>
    </>
  );
}
