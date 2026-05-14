"use client";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";
import { useDirection } from "@/shared/hooks/useDirection";
import { ChevronDown } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";
import {
  SORT_OPTIONS,
  useReservationFilters,
} from "../../_hooks/useReservationFilters";
import { useRouter } from "next/navigation";

export function SortDropdown() {
  const router = useRouter();
  const t = useTranslations("MyAccount.reservations");
  const dir = useDirection();

  const { sortBy, searchParams, lang } = useReservationFilters();

  function clickHandler(option: (typeof SORT_OPTIONS)[number]) {
    if (option === SORT_OPTIONS[0]) {
      const nextQuery = new URLSearchParams(searchParams.toString());
      nextQuery.delete("sortBy");

      const queryString = nextQuery.toString();
      const basePath = `/${lang}/reservations`;
      router.push(queryString ? `${basePath}?${queryString}` : basePath);
    } else {
      const nextQuery = new URLSearchParams(searchParams.toString());
      nextQuery.set("sortBy", option);

      const queryString = nextQuery.toString();
      router.push(`/${lang}/reservations?${queryString}`);
    }
  }

  return (
    <DropdownMenu dir={dir}>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className={cn(
            "w-50 flex items-center justify-between bg-white border rounded-lg px-4 h-12 cursor-pointer",
            "text-sm font-normal font-[inherit]",
          )}
        >
          <span className={"text-sm font-normal text-muted-foreground"}>
            {t(sortBy)}
          </span>
          <ChevronDown className="w-4 h-4 text-muted shrink-0" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        align="start"
        className="w-(--radix-dropdown-menu-trigger-width) items-start"
      >
        {SORT_OPTIONS.map((option) => (
          <DropdownMenuItem onClick={() => clickHandler(option)} key={option}>
            {t(option)}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
