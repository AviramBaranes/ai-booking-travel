import "@/app/globals.css";
import type { Metadata } from "next";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { Providers } from "../admin/Providers";
import AccountingShell from "./AccountingShell";

export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export default async function AccountingRootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const messages = await getMessages();

  return (
    <html lang="he" dir="rtl" className="h-full antialiased">
      <body className="h-full">
        <Providers>
          <NextIntlClientProvider locale="he" messages={messages}>
            <AccountingShell>{children}</AccountingShell>
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}
