import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { CheckIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const languages = [
  { code: "he", flag: "🇮🇱" },
  { code: "en", flag: "🇺🇸" },
];

export function LangSwitcher({ lang }: { lang: string }) {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const t = useTranslations("LangSwitcher");

  function handleSelect(newLang: string) {
    const rest = pathname.replace(/^\/[^/]+/, "");
    const query = searchParams.toString();
    router.push(`/${newLang}${rest}${query ? `?${query}` : ""}`);
  }

  const current = languages.find((l) => l.code === lang);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="outline">
          {/* <GlobeIcon className="size-4" /> */}
          {current?.flag} {t(lang as "he" | "en")}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {languages.map((l) => (
          <DropdownMenuItem
            key={l.code}
            onClick={() => handleSelect(l.code)}
            className="gap-2"
          >
            <span>{l.flag}</span>
            {t(l.code as "he" | "en")}
            {l.code === lang && <CheckIcon className="ms-auto size-4" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
