"use client";

import { useTranslations } from "next-intl";

interface SearchDataBannerPrintProps {
  pickUpLocationName: string;
  dropOffLocationName: string;
  pickUpDate: Date;
  dropOffDate: Date;
  pickUpTime: string;
  dropOffTime: string;
  driverAge: number;
}

function formatDateTime(date: Date, time: string) {
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}.${month}.${year} | ${time}`;
}

function formatDriverAge(age: number) {
  if (age >= 30 && age <= 65) {
    return "30 - 65";
  }
  return age;
}

export function SearchDataBannerPrint({
  pickUpLocationName,
  dropOffLocationName,
  pickUpDate,
  dropOffDate,
  pickUpTime,
  dropOffTime,
  driverAge,
}: SearchDataBannerPrintProps) {
  const t = useTranslations("booking.banner");

  return (
    <div className="hidden print:block rounded-xl border-2 border-navy p-6">
      <div className="flex justify-between gap-8">
        <div className="flex-1">
          <p className="text-sm font-semibold text-navy mb-1">{t("pickup")}</p>
          <p className="text-base font-bold text-navy">{pickUpLocationName}</p>
          <p className="text-sm text-gray-700">
            {formatDateTime(pickUpDate, pickUpTime)}
          </p>
        </div>
        <div className="w-px bg-navy/20" />
        <div className="flex-1">
          <p className="text-sm font-semibold text-navy mb-1">{t("dropoff")}</p>
          <p className="text-base font-bold text-navy">{dropOffLocationName}</p>
          <p className="text-sm text-gray-700">
            {formatDateTime(dropOffDate, dropOffTime)}
          </p>
        </div>
      </div>
      <div className="mt-3 pt-3 border-t border-navy/20">
        <p className="text-sm text-gray-700">
          {t("driverAge", { age: formatDriverAge(driverAge) })}
        </p>
      </div>
    </div>
  );
}
