import "@/app/globals.css";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { authOptions } from "@/shared/auth/authOptions";
import { Providers } from "../admin/Providers";
import AccountingShell from "./AccountingShell";

export default async function AccountingRootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect("/he/");
  }

  if (session.user?.role !== "accountant") {
    redirect("/he/");
  }

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
