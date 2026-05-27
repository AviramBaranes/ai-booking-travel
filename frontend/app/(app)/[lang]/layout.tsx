import "@/app/globals.css";
import localFont from "next/font/local";
import { getServerSession } from "next-auth/next";
import { redirect } from "next/dist/client/components/navigation";
import { notFound } from "next/navigation";
import { cookies } from "next/headers";
import { authOptions } from "@/shared/auth/authOptions";
import { Toaster } from "@/components/ui/sonner";
import { setRequestLocale } from "next-intl/server";

const polin = localFont({
  src: [
    { path: "../../fonts/Polin-Regular.otf", weight: "400" },
    { path: "../../fonts/Polin-Semibold.otf", weight: "600" },
    { path: "../../fonts/Polin-Bold.otf", weight: "700" },
    { path: "../../fonts/Polin-Black.otf", weight: "900" },
  ],
  variable: "--font-polin",
  display: "swap",
});

export default async function LangRootLayout({
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

  setRequestLocale(lang);

  return (
    <html
      lang={lang}
      dir={lang === "he" ? "rtl" : "ltr"}
      className={`h-full antialiased ${polin.variable}`}
    >
      <body className={`min-h-full flex flex-col`}>
        {children}
        <Toaster />
      </body>
    </html>
  );
}
