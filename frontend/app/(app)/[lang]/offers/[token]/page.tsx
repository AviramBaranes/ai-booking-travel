import { getClientPriceOffer } from "@/shared/api/price-offers-api";
import { ClientOfferCarCard } from "./_components/OfferCarCard";
import {
  dehydrate,
  HydrationBoundary,
  QueryClient,
} from "@tanstack/react-query";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchAddonsGallery, fetchSuppliersGallery } from "@/shared/server/cms";
import { ClientOfferSummary } from "./_components/ClientOfferSummary/ClientOfferSummary";
import { addonsGalleryKey } from "@/shared/hooks/useAddonsGallery";
import { notFound } from "next/navigation";
export default async function OfferPage({
  params,
}: {
  params: Promise<{ token: string; lang: string }>;
}) {
  const { token, lang } = await params;

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnMount: false,
        refetchOnWindowFocus: false,
      },
    },
  });
  const [offer, _, addonsGallery] = await Promise.all([
    queryClient.fetchQuery({
      queryKey: ["priceOffer", token],
      queryFn: () => getClientPriceOffer(token, () => notFound()),
    }),
    queryClient.fetchQuery({
      queryKey: suppliersGalleryKey,
      queryFn: fetchSuppliersGallery,
    }),
    queryClient.fetchQuery({
      queryKey: addonsGalleryKey,
      queryFn: fetchAddonsGallery,
    }),
  ]);

  return (
    <main className="lg:w-2/3 mx-5 lg:mx-auto lg:pt-15 pb-6">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <div className="flex flex-col-reverse lg:flex-row gap-2 mt-6">
          <div className="lg:w-3/4">
            <ClientOfferSummary
              offer={offer}
              lang={lang}
              addonsGallery={addonsGallery}
            />
          </div>
          <div className="lg:w-1/4">
            <ClientOfferCarCard offer={offer} />
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
