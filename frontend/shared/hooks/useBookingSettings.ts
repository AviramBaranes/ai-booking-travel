import { BookingSetting } from "@/payload-types";
import { useSuspenseQuery } from "@tanstack/react-query";

export const bookingSettingsKey = ["cms", "bookingSettings"] as const;

export function useBookingSettings() {
  return useSuspenseQuery<BookingSetting>({
    queryKey: bookingSettingsKey,
    queryFn: async () => {
      const res = await fetch("/api/globals/booking-settings");
      if (!res.ok) throw new Error("Failed to fetch BookingSettings");
      return res.json();
    },
  });
}
