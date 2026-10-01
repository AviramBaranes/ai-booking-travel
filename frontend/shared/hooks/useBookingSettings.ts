import { BookingSetting } from "@/payload-types";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useLocale } from "next-intl";

// The settings are localized in the CMS, so each language is cached separately.
export const bookingSettingsKey = (locale: string) =>
  ["cms", "bookingSettings", locale] as const;

export function useBookingSettings() {
  const locale = useLocale();

  return useSuspenseQuery<BookingSetting>({
    queryKey: bookingSettingsKey(locale),
    queryFn: async () => {
      const res = await fetch(
        `/api/globals/booking-settings?locale=${encodeURIComponent(locale)}`,
      );
      if (!res.ok) throw new Error("Failed to fetch BookingSettings");
      return res.json();
    },
  });
}
