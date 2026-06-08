import { getRequestConfig } from "next-intl/server";
import { notFound } from "next/navigation";

const SUPPORTED_LOCALES = ["he", "en"] as const;
type SupportedLocale = (typeof SUPPORTED_LOCALES)[number];

export default getRequestConfig(async ({ requestLocale }) => {
  const locale = (await requestLocale) ?? "he";

  if (!SUPPORTED_LOCALES.includes(locale as SupportedLocale)) {
    notFound();
  }

  return {
    locale,
    messages: (await import(`../messages/${locale}.json`)).default,
  };
});
