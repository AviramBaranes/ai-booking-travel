import { BookingStepper } from "../_components/BookingStepper";

export default async function ResultsPage() {
  await new Promise((resolve) => setTimeout(resolve, 100));
  return (
    <main className="w-2/3 mx-auto pt-10">
      <BookingStepper currentStep="results" />
    </main>
  );
}
