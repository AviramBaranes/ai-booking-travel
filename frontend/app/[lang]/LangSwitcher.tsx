"use client";

import { usePathname, useRouter } from "next/navigation";

const languages = [
  { code: "he", label: "עברית" },
  { code: "en", label: "English" },
];

export default function LangSwitcher({ lang }: { lang: string }) {
  const pathname = usePathname();
  const router = useRouter();

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const newLang = e.target.value;
    // Replace /currentLang/... with /newLang/...
    const rest = pathname.replace(/^\/[^/]+/, "");
    router.push(`/${newLang}${rest}`);
  }

  return (
    <select
      value={lang}
      onChange={handleChange}
      className="rounded border border-gray-300 bg-white px-2 py-1 text-sm"
    >
      {languages.map((l) => (
        <option key={l.code} value={l.code}>
          {l.label}
        </option>
      ))}
    </select>
  );
}
