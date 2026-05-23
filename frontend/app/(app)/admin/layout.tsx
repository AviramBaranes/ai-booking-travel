import "@/app/globals.css";
import { getServerSession } from "next-auth";
import localFont from "next/font/local";
import { redirect } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";

import { authOptions } from "@/shared/auth/authOptions";
import AdminShell from "./AdminShell";
import { Metadata } from "next";
import { Providers } from "./Providers";

export const metadata: Metadata = {
  title: "BT Admin Panel",
  description: "AI Booking Travel Admin Panel",
};

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

export default async function AdminRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect("/he/");
  }

  if (session?.user?.role !== "admin") {
    redirect("/he/");
  }

  const messages = await getMessages();

  return (
    <html lang="he" dir="rtl" className={`h-full antialiased ${polin.variable}`}>
      <body className="h-full font-sans">
        <Providers>
          <NextIntlClientProvider locale="he" messages={messages}>
            <AdminShell>{children}</AdminShell>
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}
