"use client";

import { NextIntlClientProvider } from "next-intl";
import { AuthTokenProvider } from "@/shared/components/providers/AuthTokenProvider";
import { SessionProvider } from "next-auth/react";
import { QueryProvider } from "../../_components/providers/QueryProvider";

interface BookingProvidersProps {
  children: React.ReactNode;
  lang: string;
  messages: Record<string, string>;
}
export function BookingProviders({
  children,
  lang,
  messages,
}: BookingProvidersProps) {
  return (
    <QueryProvider>
      <NextIntlClientProvider locale={lang} messages={messages}>
        <SessionProvider>
          <AuthTokenProvider>{children}</AuthTokenProvider>
        </SessionProvider>
      </NextIntlClientProvider>
    </QueryProvider>
  );
}
