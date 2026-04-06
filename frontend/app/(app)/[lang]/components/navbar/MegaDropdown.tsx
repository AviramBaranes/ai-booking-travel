"use client";

import Image from "next/image";
import Link from "next/link";
import { useState, useRef } from "react";
import type { Page } from "@/payload-types";
import type { Populated } from "@/shared/types/payload";

interface MegaDropdownProps {
  label: string;
  links: { label: string; page: number | Page; id?: string | null }[];
  lang: string;
}

export function MegaDropdown({ label, links, lang }: MegaDropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleMouseEnter = () => {
    if (timeoutRef.current) clearTimeout(timeoutRef.current);
    setIsOpen(true);
  };

  const handleMouseLeave = () => {
    timeoutRef.current = setTimeout(() => setIsOpen(false), 150);
  };

  return (
    <div
      className="relative"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <button
        className="flex cursor-pointer items-center gap-1.5 type-h6 text-navy"
        onClick={() => setIsOpen((v) => !v)}
      >
        {label}
        <Image
          src="/assets/header/chevron-down.svg"
          alt=""
          width={8}
          height={7}
          className={`transition-transform ${isOpen ? "rotate-180" : ""}`}
        />
      </button>

      <div
        className={`absolute inset-s-0 top-full z-50 mt-4 min-w-56 rounded-xl bg-white py-2 shadow-[0px_4px_16px_rgba(15,0,67,0.15)] ${isOpen ? "" : "hidden"}`}
      >
        {links.map((link, i) => {
          const page = link.page as Populated<typeof link.page>;
          return (
            <div key={link.id}>
              {i > 0 && <div className="mx-3 border-t border-border-light" />}
              <Link
                href={`/${lang}/${page.slug}`}
                className="block px-4 py-2.5 text-base font-semibold text-navy hover:bg-brand-blue/10 hover:text-brand-blue"
                onClick={() => setIsOpen(false)}
              >
                {link.label}
              </Link>
            </div>
          );
        })}
      </div>
    </div>
  );
}
