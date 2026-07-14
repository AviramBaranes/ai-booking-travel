import { Suspense } from "react";
import { HydrationBoundary, dehydrate } from "@tanstack/react-query";
import { BookingStepper } from "../_components/BookingStepper";
import { OrderPageContent } from "./_components/OrderPageContent";
import { OrderLoadError } from "./_components/OrderLoadError";
import { ExpiredSearchGate } from "../_components/ExpiredSearchGate";
import { BackButton } from "@/shared/components/booking/BackButton";
import { SearchDataBannerWithQuery } from "../_components/SearchDataBannerWithQuery";
import { QueryErrorBoundary } from "@/shared/components/QueryErrorBoundary";
import { Loading } from "@/shared/components/Loading";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { bookingSettingsKey } from "@/shared/hooks/useBookingSettings";
import { fetchBookingSettings } from "@/shared/server/cms";

export default async function OrderPage() {
  const queryClient = getQueryClient();

  await queryClient.prefetchQuery({
    queryKey: bookingSettingsKey,
    queryFn: fetchBookingSettings,
  });

  return (
    <main className="lg:w-2/3 mx-auto lg:pt-15 pb-6">
      <BookingStepper currentStep="ordering" />
      <div className="my-4">
        <SearchDataBannerWithQuery />
      </div>
      <BackButton translationKey="backToPlans" />
      <ExpiredSearchGate>
        <QueryErrorBoundary fallback={<OrderLoadError />}>
          <Suspense fallback={<Loading className="mx-auto mt-10" />}>
            <HydrationBoundary state={dehydrate(queryClient)}>
              <OrderPageContent />
            </HydrationBoundary>
          </Suspense>
        </QueryErrorBoundary>
      </ExpiredSearchGate>
    </main>
  );
}
