import { Header } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { ChevronDown } from "lucide-react";
import Link from "next/link";
import { useState } from "react";

type HeaderLink = NonNullable<Header["links"]>[number];

export function MobileMegaRow({
  link,
  lang,
  onNavigate,
}: {
  link: HeaderLink;
  lang: string;
  onNavigate: () => void;
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="w-full">
      <button
        className="flex h-13 w-full items-center px-4"
        onClick={() => setExpanded((v) => !v)}
      >
        <span className="type-h6 text-navy">{link.megaLabel}</span>
        <ChevronDown
          className={`ms-auto size-4 text-gray-400 transition-transform ${expanded ? "rotate-180" : ""}`}
        />
      </button>
      <div className="bg-border-light h-px w-full" />

      {expanded &&
        link.megaLinks?.map((sub) => {
          const page = sub.page as Populated<typeof sub.page>;
          return (
            <div key={sub.id}>
              <Link
                href={`/${lang}/${page.slug}`}
                className="flex h-12 w-full items-center px-6 type-h6 text-navy/70"
                onClick={onNavigate}
              >
                {sub.label}
              </Link>
              <div className="bg-border-light h-px w-full" />
            </div>
          );
        })}
    </div>
  );
}
