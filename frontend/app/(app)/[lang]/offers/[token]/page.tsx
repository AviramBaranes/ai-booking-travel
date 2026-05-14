import { getClientPriceOffer } from "@/shared/api/price-offers-api";
export default async function OfferPage({
  params,
}: {
  params: Promise<{ token: string }>;
}) {
  const { token } = await params;
  const offer = await getClientPriceOffer(token);

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <h1 className="text-2xl font-bold mb-4">{offer.name}</h1>
      {Object.entries(offer).map(([key, value]) => (
        <div key={key} className="mb-2">
          <span className="font-semibold">{key}: </span>
          <span>
            {typeof value === "object" ? JSON.stringify(value) : value}
          </span>
        </div>
      ))}
    </main>
  );
}
