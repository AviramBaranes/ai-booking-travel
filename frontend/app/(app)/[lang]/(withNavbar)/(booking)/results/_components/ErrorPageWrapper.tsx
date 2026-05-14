"use client";

import { SessionProvider } from "next-auth/react";
import { NextIntlClientProvider } from "next-intl";
import ErrorResultPageContent from "./ErrorPage";

export function ErrorPageWrapper({
  locale,
  messages,
}: {
  locale: string;
  messages: Record<string, unknown>;
}) {
  return (
    <NextIntlClientProvider locale={locale} messages={messages}>
      <SessionProvider>
        <ErrorResultPageContent />
      </SessionProvider>
    </NextIntlClientProvider>
  );
}
