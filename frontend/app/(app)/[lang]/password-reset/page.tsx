import type { Metadata } from "next";
import { SUPPORTED_LANGS } from "@/shared/constants/supported_langs";
import Link from "next/link";
import Image from "next/image";
import { PasswordResetForm } from "./_components/PasswordResetForm";
import { NextIntlClientProvider } from "next-intl";
import { getMessages } from "next-intl/server";
import { AppProviders } from "../_components/providers/AppProviders";
import { Suspense } from "react";

export const metadata: Metadata = {
  robots: { index: false, follow: false },
};

export async function generateStaticParams() {
  const params = SUPPORTED_LANGS.map((locale) => ({ lang: locale }));
  return params;
}

export default async function PasswordResetPage({
  params,
}: {
  params: Promise<{ lang: string }>;
}) {
  const { lang } = await params;
  const messages = await getMessages();

  return (
    <main className="min-h-screen flex flex-col items-center justify-center bg-background p-4">
      <div className="w-full max-w-100 flex flex-col items-center gap-8">
        <Link href={`/${lang}`}>
          <Image
            src="/logo.png"
            alt="AIBookingTravel"
            width={168}
            height={32}
            className="object-contain"
            priority
          />
        </Link>
        <AppProviders lang={lang} messages={messages}>
          <Suspense fallback={null}>
          <PasswordResetForm />
          </Suspense>
        </AppProviders>
      </div>
    </main>
  );
}
