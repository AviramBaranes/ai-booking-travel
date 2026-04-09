"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SessionProvider } from "next-auth/react";
import { NextIntlClientProvider } from "next-intl";

export function AppProviders({
  children,
  lang,
  messages,
}: {
  children: React.ReactNode;
  lang: string;
  messages?: Record<string, unknown>;
}) {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <SessionProvider>
        <NextIntlClientProvider locale={lang} messages={messages}>
          {children}
        </NextIntlClientProvider>
      </SessionProvider>
    </QueryClientProvider>
  );
}
