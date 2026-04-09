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
  const availability = await queryClient.fetchQuery({
    queryKey: bookingKeys.availability(searchRequest),
    queryFn: () => searchAvailableCars(searchRequest),
  });

  return (
    <main className="w-2/3 mx-auto pt-15 pb-300">
      <BookingStepper currentStep="results" />
      <NextIntlClientProvider locale={lang} messages={messages}>
        <div className="my-4">
          <SearchDataBanner
            pickUpLocation={{
              id: query.pickupLocationId,
              name: availability.pickupLocationName,
            }}
            dropOffLocation={{
              id: query.returnLocationId,
              name: availability.dropoffLocationName,
            }}
            pickUpTime={query.pickupTime}
            dropOffTime={query.returnTime}
            pickUpDate={query.pickupDate}
            dropOffDate={query.returnDate}
            driverAge={query.driverAge}
            couponCode={query.couponCode}
            showButton
          />
        </div>
        <HydrationBoundary state={dehydrate(queryClient)}>
          <CarResults searchRequest={searchRequest} />
        </HydrationBoundary>
      </NextIntlClientProvider>
    </main>
  );
}
