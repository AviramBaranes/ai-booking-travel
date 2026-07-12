import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { PriceOfferCarCard } from "./_components/PriceOfferCarCard";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import {
  fetchAddonsGallery,
  fetchBookingSettings,
  fetchSuppliersGallery,
} from "@/shared/server/cms";
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { BackButton } from "@/shared/components/booking/BackButton";
import { PriceOfferSummary } from "./_components/PriceOfferSummary/PriceOfferSummary";
import { bookingSettingsKey } from "@/shared/hooks/useBookingSettings";
import { addonsGalleryKey } from "@/shared/hooks/useAddonsGallery";

export default async function PriceOfferPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const lang = await getLang();
  const { id } = await params;

  if (!id || isNaN(Number(id))) {
    redirect(`/${lang}/price-offers`);
  }

  const queryClient = getQueryClient();
  await Promise.all([
    queryClient.fetchQuery({
      queryKey: suppliersGalleryKey,
      queryFn: fetchSuppliersGallery,
    }),
    queryClient.fetchQuery({
      queryKey: bookingSettingsKey,
      queryFn: fetchBookingSettings,
    }),
    queryClient.fetchQuery({
      queryKey: addonsGalleryKey,
      queryFn: fetchAddonsGallery,
    }),
  ]);

  return (
    <main className="w-2/3 mx-auto pt-4 pb-6">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <BackButton
          translationKey="backToPriceOffers"
          href={`/${lang}/price-offers`}
        />
        <div className="flex gap-2 mt-6">
          <div className="w-3/4">
            <PriceOfferSummary priceOfferId={Number(id)} />
          </div>
          <div className="w-1/4">
            <PriceOfferCarCard priceOfferId={Number(id)} />
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
