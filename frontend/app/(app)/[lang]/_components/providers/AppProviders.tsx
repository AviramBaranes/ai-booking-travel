"use client";

import { AuthTokenProvider } from "@/shared/components/providers/AuthTokenProvider";
import { SessionProvider } from "next-auth/react";
import { NextIntlClientProvider } from "next-intl";
import { QueryProvider } from "./QueryProvider";

export function AppProviders({
  children,
  lang,
  messages,
  showDevtools = false,
}: {
  children: React.ReactNode;
  lang: string;
  messages?: Record<string, unknown>;
  showDevtools?: boolean;
}) {
  return (
    <QueryProvider showDevtools={showDevtools}>
      <NextIntlClientProvider locale={lang} messages={messages}>
        <SessionProvider>
          <AuthTokenProvider>{children}</AuthTokenProvider>
        </SessionProvider>
      </NextIntlClientProvider>
    </QueryProvider>
  );
}
