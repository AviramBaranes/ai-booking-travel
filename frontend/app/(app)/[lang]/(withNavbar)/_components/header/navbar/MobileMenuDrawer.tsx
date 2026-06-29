"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, X } from "lucide-react";

import {
  Sheet,
  SheetContent,
  SheetClose,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Logo } from "./Logo";
import type { Header } from "@/payload-types";
import type { Populated } from "@/shared/types/payload";
import type { Page } from "@/payload-types";
import { MobileMegaRow } from "./MobileMegaRow";

interface Props {
  lang: string;
  links: NonNullable<Header["links"]>;
}

export function MobileMenuDrawer({ lang, links }: Props) {
  const [open, setOpen] = useState(false);

  const close = () => setOpen(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <button aria-label="Open menu" className="lg:hidden">
          <Menu className="size-6" />
        </button>
      </SheetTrigger>

      <SheetContent
        side="top"
        showCloseButton={false}
        className="p-0 rounded-none border-0 w-full"
      >
        <SheetTitle className="sr-only">Navigation menu</SheetTitle>
        {/* Top bar */}
        <div className="flex h-15 items-center justify-between px-4 border-b border-border-light">
          <SheetClose asChild>
            <button aria-label="Close menu">
              <X className="size-5 text-navy" />
            </button>
          </SheetClose>
          <Logo lang={lang} />
        </div>

        {/* Nav rows */}
        <nav aria-label="Mobile navigation" className="flex flex-col">
          {links.map((link) =>
            link.type === "mega" ? (
              <MobileMegaRow
                key={link.id}
                link={link}
                lang={lang}
                onNavigate={close}
              />
            ) : (
              <div key={link.id} className="w-full">
                <Link
                  href={`/${lang}/${(link.page as Populated<Page>)?.slug ?? ""}`}
                  className="flex h-13 w-full items-center px-4 type-h6 text-navy"
                  onClick={close}
                >
                  {link.label}
                </Link>
                <div className="bg-border-light h-px w-full" />
              </div>
            ),
          )}
        </nav>
        <div className="h-14"/>
      </SheetContent>
    </Sheet>
  );
}
