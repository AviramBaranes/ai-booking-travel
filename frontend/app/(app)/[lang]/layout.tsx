import "@/app/globals.css";
import { getServerSession } from "next-auth/next";
import Providers from "../providers";
import { redirect } from "next/dist/client/components/navigation";
import { notFound } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { authOptions } from "@/shared/auth/authOptions";
import { Navbar } from "./components/navbar/Navbar";
import { Footer } from "./components/footer/Footer";
import { BackToAdminBanner } from "./components/BackToAdminBanner";

export default async function AppRootLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const session = await getServerSession(authOptions);
  const isAdmin = session?.user?.role === "admin";
  if (isAdmin) {
    redirect("/admin/");
  }
  const { lang } = await params;

  if (!["he", "en"].includes(lang)) {
    notFound();
  }

  return (
    <html
      lang={lang}
      dir={lang === "he" || lang === "ar" ? "rtl" : "ltr"}
      className="h-full antialiased"
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          <NextIntlClientProvider locale={lang}>
            <BackToAdminBanner />
            <Navbar lang={lang} isAuthenticated={!!session?.user?.id} />
            {children}
            <Footer lang={lang} />
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}
