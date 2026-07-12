"use client";

import { AuthBootstrap } from "@/shared/components/providers/AuthBootstrap";
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
        <AuthBootstrap>{children}</AuthBootstrap>
      </NextIntlClientProvider>
    </QueryProvider>
  );
}
