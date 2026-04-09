import { SearchDataBanner } from "@/shared/components/booking/SearchDataBanner";
import { BookingStepper } from "../_components/BookingStepper";

export default async function ResultsPage() {
  await new Promise((resolve) => setTimeout(resolve, 100));
  return (
    <main className="w-2/3 mx-auto pt-10 pb-300">
      <BookingStepper currentStep="results" />
      <div className="my-4">
        <SearchDataBanner
          pickUpLocation={{
            id: 1,
            name: "Holland, Amsterdam - Schiphol Airport",
          }}
          dropOffLocation={{
            id: 1,
            name: "Holland, Amsterdam - Schiphol Airport",
          }}
          pickUpTime="12:00"
          dropOffTime="10:00"
          pickUpDate={new Date(2025, 4, 10)}
          dropOffDate={new Date(2025, 4, 16)}
          driverAge={28}
          couponCode="test"
          showButton
        />
      </div>
    </main>
  );
}
