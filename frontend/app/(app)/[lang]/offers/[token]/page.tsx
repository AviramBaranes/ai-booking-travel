import { getClientPriceOffer } from "@/shared/api/price-offers-api";
import { ClientOfferCarCard } from "./_components/OfferCarCard";
import {
  dehydrate,
  HydrationBoundary,
  QueryClient,
} from "@tanstack/react-query";
import { suppliersGalleryKey } from "@/shared/hooks/useSuppliersGallery";
import { fetchSuppliersGallery } from "@/shared/server/cms";
import { ClientOfferSummary } from "./_components/ClientOfferSummary/ClientOfferSummary";
export default async function OfferPage({
  params,
}: {
  params: Promise<{ token: string, lang: string }>;
}) {
  const { token, lang } = await params;
  const offer = await getClientPriceOffer(token);

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnMount: false,
        refetchOnWindowFocus: false,
      },
    },
  });
  await queryClient.fetchQuery({
    queryKey: suppliersGalleryKey,
    queryFn: fetchSuppliersGallery,
  });

  return (
    <main className="w-2/3 mx-auto pt-4 pb-6">
      <HydrationBoundary state={dehydrate(queryClient)}>
        <div className="flex gap-2 mt-6">
          <div className="w-3/4">
          <ClientOfferSummary offer={offer} lang={lang} />
          </div>
          <div className="w-1/4">
            <ClientOfferCarCard offer={offer} />
          </div>
        </div>
      </HydrationBoundary>
    </main>
  );
}
