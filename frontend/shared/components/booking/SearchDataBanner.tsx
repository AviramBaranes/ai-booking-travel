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
import { booking } from "@/shared/client";
import { useAvailableCars } from "@/shared/hooks/useAvailableCars";

interface SearchDataBannerProps extends Omit<
  SearchFormFields,
  "pickUpLocation" | "dropOffLocation"
> {
  pickUpLocationId: number;
  dropOffLocationId: number;
  searchRequest: booking.SearchAvailabilityRequest;
  showButton?: boolean;
  fromCache?: boolean;
}

function formatDriverAge(age: number) {
  if (age >= 30 && age <= 65) {
    return "30 - 65";
  }

  return age;
}

export function SearchDataBanner({
  pickUpLocationId,
  dropOffLocationId,
  pickUpTime,
  dropOffTime,
  pickUpDate,
  dropOffDate,
  driverAge,
  couponCode,
  showButton,
  searchRequest,
  fromCache,
}: SearchDataBannerProps) {
  const t = useTranslations("booking.banner");
  const dir = useDirection();
  const { data } = useAvailableCars(searchRequest, { fromCache });
  const pickUpLocationName = data?.pickupLocationName ?? "";
  const dropOffLocationName = data?.dropoffLocationName ?? "";
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
          pickUpLocation={{
            id: pickUpLocationId,
            name: pickUpLocationName,
          }}
          dropOffLocation={{
            id: dropOffLocationId,
            name: dropOffLocationName,
          }}
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
            location={pickUpLocationName}
            date={pickUpDate}
            time={pickUpTime}
          />

          <div className="h-25 w-px bg-white" />

          <SearchDataPoint
            icon={ArrowDown}
            label={t("dropoff")}
            location={dropOffLocationName ?? pickUpLocationName}
            date={dropOffDate}
            time={dropOffTime}
          />
        </div>
      </div>

      {showButton && (
        <Button
          variant="ghost"
          className={clsx(
            "absolute cursor-pointer -top-px flex items-center gap-2 rounded-none bg-white px-3 py-6 hover:bg-white/90",
            {
              "-left-px rounded-br-xl rounded-tl-3xl": dir === "rtl",
              "-right-px rounded-bl-xl rounded-tr-3xl": dir === "ltr",
            },
          )}
          onClick={() => setShowForm(true)}
        >
          <Pencil className="size-4 text-brand" />
          <span className="text-[17px] font-semibold text-brand">
            {t("editSearch")}
          </span>
        </Button>
      )}

      <div
        className={clsx(
          "absolute bottom-0 flex items-center gap-2 border-t border-white px-3 py-4",
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
