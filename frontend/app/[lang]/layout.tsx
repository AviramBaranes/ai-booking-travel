import { getServerSession } from "next-auth/next";
import "../globals.css";
import Providers from "../providers";
import { LangSwitcher } from "./LangSwitcher";
import { redirect } from "next/dist/client/components/navigation";
import { LoginModal } from "./LoginModal";
import { NextIntlClientProvider } from "next-intl";
import { authOptions } from "@/shared/auth/authOptions";

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
  return (
    <html
      lang={lang}
      dir={lang === "he" || lang === "ar" ? "rtl" : "ltr"}
      className="h-full antialiased"
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          <NextIntlClientProvider locale={lang}>
            <header className="flex items-center justify-between px-4 py-2 border-b">
              <span className="font-semibold">AI Booking Travel</span>
              <LangSwitcher lang={lang} />
              {!session?.user?.id && <LoginModal />}
            </header>
            {children}
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}
