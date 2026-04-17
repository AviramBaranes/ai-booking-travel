"use client";

import { User, ArrowUp, ArrowDown, Pencil } from "lucide-react";
import { SearchDataPoint } from "./SearchDataPoint";
import { Button } from "@/components/ui/button";
import { clsx } from "clsx";
import { useTranslations } from "next-intl";
import { useDirection } from "@/shared/hooks/useDirection";

function formatDriverAge(age: number) {
  if (age >= 30 && age <= 65) {
    return "30 - 65";
  }
  return age;
}

export interface SearchDataBannerDisplayProps {
  pickUpLocationName: string;
  dropOffLocationName: string;
  pickUpDate: Date;
  dropOffDate: Date;
  pickUpTime: string;
  dropOffTime: string;
  driverAge: number;
  showButton?: boolean;
  onEditClick?: () => void;
}

export function SearchDataBannerDisplay({
  pickUpLocationName,
  dropOffLocationName,
  pickUpDate,
  dropOffDate,
  pickUpTime,
  dropOffTime,
  driverAge,
  showButton,
  onEditClick,
}: SearchDataBannerDisplayProps) {
  const t = useTranslations("booking.banner");
  const dir = useDirection();

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
            location={dropOffLocationName}
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
          onClick={onEditClick}
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
