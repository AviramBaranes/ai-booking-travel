import type { Metadata } from "next";
import "../globals.css";
import Providers from "../providers";
import LangSwitcher from "./LangSwitcher";

export const metadata: Metadata = {
  title: "Home",
  description: "AI Booking Travel - Find your next car to rent worldwide",
};

export default async function LangLayout({
  children,
  params,
}: Readonly<{
  children: React.ReactNode;
  params: Promise<{ lang: string }>;
}>) {
  const { lang } = await params;
  return (
    <html
      lang={lang}
      dir={lang === "he" || lang === "ar" ? "rtl" : "ltr"}
      className="h-full antialiased"
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          <header className="flex items-center justify-between px-4 py-2 border-b">
            <span className="font-semibold">AI Booking Travel</span>
            <LangSwitcher lang={lang} />
          </header>
          {children}
        </Providers>
      </body>
    </html>
  );
}
