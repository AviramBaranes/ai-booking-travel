"use client";

import { usePathname, useRouter } from "next/navigation";
import { useTranslations } from "next-intl";

const languages = [
  { code: "he", flag: "🇮🇱" },
  { code: "en", flag: "🇺🇸" },
];

export function LangSwitcher({ lang }: { lang: string }) {
  const pathname = usePathname();
  const router = useRouter();
  const t = useTranslations("LangSwitcher");

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const newLang = e.target.value;
    const rest = pathname.replace(/^\/[^/]+/, "");
    router.push(`/${newLang}${rest}`);
  }

  return (
    <select
      value={lang}
      onChange={handleChange}
      className="cursor-pointer appearance-none rounded-full border-2 border-navy bg-white px-4 py-2 text-base font-bold text-navy"
    >
      {languages.map((l) => (
        <option key={l.code} value={l.code}>
          {l.flag} {t(l.code as "he" | "en")}
        </option>
      ))}
    </select>
  );
}
