import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BookingStepper } from "../_components/BookingStepper";
import { getMessages } from "next-intl/server";
import { getLang } from "@/shared/lang/lang";
import { NextIntlClientProvider } from "next-intl";
import { redirect } from "next/navigation";
import { parseSearchQuery, toSearchRequest } from "./searchQuery";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { bookingKeys } from "@/shared/hooks/useAvailableCars";
import { searchAvailableCars } from "@/shared/api/booking-api";
import { CarResults } from "./CarResults";
import { ErrorPageWrapper } from "./_components/ErrorPageWrapper";
import {
  fetchSuppliersGallery,
  fetchAddonsGallery,
  fetchBookingSettings,
} from "@/shared/server/cms";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { addonsGalleryKey } from "@/shared/hooks/useAddonsGallery";
import { bookingSettingsKey } from "@/shared/hooks/useBookingSettings";

export default async function ResultsPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string>>;
}) {
  const lang = await getLang();
  const resolvedParams = await searchParams;
  const query = parseSearchQuery(new URLSearchParams(resolvedParams));

  if (!query) {
    redirect(`/${lang}`);
  }
  const messages = await getMessages({ locale: lang });

  const searchRequest = toSearchRequest(query);
  const queryClient = getQueryClient();

  try {
    const [result] = await Promise.all([
      queryClient.fetchQuery({
        queryKey: bookingKeys.availability(searchRequest),
        queryFn: () => searchAvailableCars(searchRequest),
      }),
      queryClient.fetchQuery({
        queryKey: suppliersGalleryKey,
        queryFn: fetchSuppliersGallery,
      }),
      queryClient.fetchQuery({
        queryKey: addonsGalleryKey,
        queryFn: fetchAddonsGallery,
      }),
      queryClient.fetchQuery({
        queryKey: bookingSettingsKey,
        queryFn: fetchBookingSettings,
      }),
    ]);
    if (!result.availableVehicles.length) throw new Error("No results");
  } catch {
    return <ErrorPageWrapper locale={lang} messages={messages} />;
  }

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <NextIntlClientProvider locale={lang} messages={messages}>
        <BookingStepper currentStep="results" />
        <HydrationBoundary state={dehydrate(queryClient)}>
          <div className="my-4">
            <SearchDataBanner
              pickUpLocationId={query.pickupLocationId}
              dropOffLocationId={query.returnLocationId}
              pickUpTime={query.pickupTime}
              dropOffTime={query.returnTime}
              pickUpDate={query.pickupDate}
              dropOffDate={query.returnDate}
              driverAge={query.driverAge}
              couponCode={query.couponCode}
              searchRequest={searchRequest}
              showButton
            />
          </div>
          <CarResults searchRequest={searchRequest} />
        </HydrationBoundary>
      </NextIntlClientProvider>
    </main>
  );
}
