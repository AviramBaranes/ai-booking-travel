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
import { cache } from "react";
import type { Metadata } from "next";
import { formatRangeDate } from "@/shared/utils/formatDate";

type Props = {
  params: Promise<{ token: string; lang: string }>;
};

const getOffer = cache((token: string) =>
  getClientPriceOffer(token, () => notFound()),
);

// For the WhatsApp link preview, not search. Excludes offer.name (usually the
// customer's) and the price: previews are cached and re-rendered wherever the
// link is forwarded.
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { token, lang } = await params;
  const offer = await getOffer(token);

  const car = offer.carDetails;
  const carName = car.model || car.carGroup;
  const dates = `${formatRangeDate(lang, new Date(offer.pickupDate))} – ${formatRangeDate(lang, new Date(offer.dropoffDate))}`;
  const days =
    lang === "he"
      ? `${offer.rentalDays} ימי השכרה`
      : `${offer.rentalDays} rental days`;

  const title = `${carName} · ${offer.pickupLocationName}`;
  const description = `${dates} · ${days}`;

  return {
    title,
    description,
    robots: { index: false, follow: false },
    openGraph: {
      title,
      description,
      type: "website",
      images: car.imageUrl ? [{ url: car.imageUrl }] : undefined,
    },
  };
}

export default async function OfferPage({ params }: Props) {
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
      queryFn: () => getOffer(token),
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
