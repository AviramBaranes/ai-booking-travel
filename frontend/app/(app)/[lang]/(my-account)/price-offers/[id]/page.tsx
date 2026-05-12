import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";
import { Suspense } from "react";
import { PriceOfferCarCard } from "./_components/PriceOfferCarCard";
import { SelectedCarCardSkeleton } from "@/shared/components/booking/SelectedCarCard/SelectedCarCardSkeleton";
import { getQueryClient } from "@/shared/hooks/getQueryClient";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchSuppliersGallery } from "@/shared/server/cms";
import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { BackButton } from "../../../../../../shared/components/booking/BackButton";
import { SummarySkeleton } from "../../_components/SummarySkeleton";
import { PriceOfferSummary } from "./_components/PriceOfferSummary/PriceOfferSummary";

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
  await queryClient.fetchQuery({
    queryKey: suppliersGalleryKey,
    queryFn: fetchSuppliersGallery,
  });

  return (
    <main className="w-2/3 mx-auto pt-4 pb-6 print:w-full">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <div className="print:hidden">
          <BackButton
            translationKey="backToPriceOffers"
            href={`/${lang}/price-offers`}
          />
        </div>
        <div className="flex gap-2 mt-6 print:flex-col print:gap-6">
          <div className="w-3/4 print:w-full">
            <Suspense fallback={<SummarySkeleton />}>
              <PriceOfferSummary priceOfferId={Number(id)} />
            </Suspense>
          </div>
          <div className="w-1/4 print:w-full">
            <Suspense fallback={<SelectedCarCardSkeleton />}>
              <PriceOfferCarCard priceOfferId={Number(id)} />
            </Suspense>
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
