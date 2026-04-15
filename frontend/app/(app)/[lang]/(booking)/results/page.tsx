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
import { getPayload } from "payload";
import config from "@payload-config";
import { ErrorPageWrapper } from "./_components/ErrorPageWrapper";

export async function getSupplierGallery() {
  const payload = await getPayload({ config });
  return payload.findGlobal({
    slug: "suppliersGallery",
    draft: false,
  });
}

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
  let supplierGallery:
    | Awaited<ReturnType<typeof getSupplierGallery>>
    | undefined;
  try {
    const [result, gallery] = await Promise.all([
      queryClient.fetchQuery({
        queryKey: bookingKeys.availability(searchRequest),
        queryFn: () => searchAvailableCars(searchRequest),
      }),
      getSupplierGallery(),
    ]);
    if (!result.availableVehicles.length) throw new Error("No results");
    supplierGallery = gallery;
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
          <CarResults
            searchRequest={searchRequest}
            supplierGallery={supplierGallery}
          />
        </HydrationBoundary>
      </NextIntlClientProvider>
    </main>
  );
}
