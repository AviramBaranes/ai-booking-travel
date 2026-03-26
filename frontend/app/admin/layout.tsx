import "../globals.css";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";

import { authOptions } from "@/shared/auth/authOptions";
import Providers from "@/app/providers";
import AdminShell from "./AdminShell";

export default async function AdminRootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await getServerSession(authOptions);
  if (!session) {
    redirect("/he/");
  }

  const isAdmin = session?.user?.role === "admin";
  if (!isAdmin) {
    redirect("/he/");
  }

  const messages = await getMessages();

  return (
    <html lang="he" dir="rtl" className="h-full antialiased">
      <body className="h-full">
        <Providers>
          <NextIntlClientProvider locale="he" messages={messages}>
            <AdminShell>{children}</AdminShell>
          </NextIntlClientProvider>
        </Providers>
      </body>
    </html>
  );
}
