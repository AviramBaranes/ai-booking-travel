"use client";

import { User, ArrowUp, ArrowDown, Pencil, X } from "lucide-react";
import { SearchDataPoint } from "./SearchDataPoint";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import {
  SearchForm,
  SearchFormFields,
} from "@/app/(app)/[lang]/_components/home/SearchForm/SearchForm";
import { useDirection } from "@/shared/hooks/useDirection";
import { clsx } from "clsx";
import { useTranslations } from "next-intl";

interface SearchDataBannerProps extends SearchFormFields {
  showButton?: boolean;
}

function formatDriverAge(age: number) {
  if (age >= 30 && age <= 65) {
    return "30 - 65";
  }

  return age;
}

export function SearchDataBanner({
  pickUpLocation,
  dropOffLocation,
  pickUpTime,
  dropOffTime,
  pickUpDate,
  dropOffDate,
  driverAge,
  couponCode,
  showButton,
}: SearchDataBannerProps) {
  const t = useTranslations("booking.banner");
  const dir = useDirection();
  const [showForm, setShowForm] = useState(false);

  if (showForm) {
    return (
      <div className="relative bg-navy py-4 px-2 rounded-xl">
        <X
          className={clsx("absolute top-2 cursor-pointer text-muted", {
            "left-2": dir === "rtl",
            "right-2": dir === "ltr",
          })}
          onClick={() => setShowForm(false)}
        />
        <SearchForm
          className="w-full"
          pickUpLocation={pickUpLocation}
          dropOffLocation={dropOffLocation}
          pickUpDate={pickUpDate}
          dropOffDate={dropOffDate}
          pickUpTime={pickUpTime}
          dropOffTime={dropOffTime}
          couponCode={couponCode}
          driverAge={driverAge}
        />
      </div>
    );
  }

  return (
    <section
      className="relative border-none w-full rounded-3xl bg-navy bg-cover bg-center bg-no-repeat shadow-card"
      style={{
        backgroundImage: "url('/assets/booking/search-data-bg.png')",
      }}
    >
      <div className="flex items-center justify-between px-10 py-3">
        <div className="flex items-center gap-13">
          <SearchDataPoint
            icon={ArrowUp}
            label={t("pickup")}
            location={pickUpLocation.name}
            date={pickUpDate}
            time={pickUpTime}
          />

          <div className="h-25 w-px bg-white" />

          <SearchDataPoint
            icon={ArrowDown}
            label={t("dropoff")}
            location={dropOffLocation?.name ?? pickUpLocation.name}
            date={dropOffDate}
            time={dropOffTime}
          />
        </div>
      </div>

      {showButton && (
        <Button
          variant="ghost"
          className={clsx(
            "absolute cursor-pointer -top-px flex items-center gap-2 rounded-none bg-white px-5 py-6.5 hover:bg-white/90",
            {
              "-left-px rounded-br-xl rounded-tl-3xl": dir === "rtl",
              "-right-px rounded-bl-xl rounded-tr-3xl": dir === "ltr",
            },
          )}
          onClick={() => setShowForm(true)}
        >
          <Pencil className="size-4 text-brand" />
          <span className="type-h6 text-brand">{t("editSearch")}</span>
        </Button>
      )}

      <div
        className={clsx(
          "absolute bottom-0 flex items-center gap-2 border-t border-white px-4 py-5",
          {
            "left-0 rounded-bl-3xl rounded-tr-3xl border-r": dir === "rtl",
            "right-0 rounded-br-3xl rounded-tl-3xl border-l": dir === "ltr",
          },
        )}
      >
        <User className="size-4 text-white" />
        <p className="type-paragraph text-center font-bold text-white whitespace-nowrap">
          {t("driverAge", { age: formatDriverAge(driverAge) })}
        </p>
      </div>
    </section>
  );
}
