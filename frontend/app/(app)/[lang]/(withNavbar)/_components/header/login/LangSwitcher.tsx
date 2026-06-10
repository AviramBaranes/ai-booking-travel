import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { CheckIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useDirection } from "@/shared/hooks/useDirection";

const languages = [
  { code: "he", flag: "🇮🇱" },
  { code: "en", flag: "🇺🇸" },
];

export function LangSwitcher({ lang }: { lang: string }) {
  const router = useRouter();
  const t = useTranslations("LangSwitcher");
  const dir = useDirection();

  function handleSelect(newLang: string) {
    const href = location.href;
    const newPathname = href.replace(`/${lang}`, `/${newLang}`);
    router.push(newPathname);
  }

  const current = languages.find((l) => l.code === lang);

  return (
    <DropdownMenu dir={dir}>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="outline">
          {/* <GlobeIcon className="size-4" /> */}
          {current?.flag} {t(lang as "he" | "en")}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        {languages.map((l) => (
          <DropdownMenuItem
            key={l.code}
            onClick={() => handleSelect(l.code)}
            className="gap-2 flex"
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
