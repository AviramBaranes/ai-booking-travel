import { getLang } from "@/shared/lang/lang";
import { redirect } from "next/dist/client/components/navigation";

export default async function ReservationDetailsPage({
  params,
}: {
  params: { reservationId: string };
}) {
  const lang = await getLang();
  const reservationId = params.reservationId;

  if (!reservationId) {
    redirect(`/${lang}/my-account/reservations`);
  }

  return (
    <main className="w-2/3 mx-auto pt-15 pb-6">
      <h1 className="text-2xl font-bold mb-4">Reservation Details</h1>
      <p>Details for reservation ID: {reservationId}</p>
    </main>
  );
}
