import "@/app/globals.css";
import { getServerSession } from "next-auth";
import { redirect } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";

import { authOptions } from "@/shared/auth/authOptions";
import AdminShell from "./AdminShell";
import { Metadata } from "next";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SessionProvider } from "next-auth/react";
import { AuthTokenProvider } from "@/shared/components/providers/AuthTokenProvider";
import { Providers } from "./Providers";

export const metadata: Metadata = {
  title: "BT Admin Panel",
  description: "AI Booking Travel Admin Panel",
};

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
