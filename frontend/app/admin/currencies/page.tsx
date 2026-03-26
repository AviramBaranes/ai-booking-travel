import CurrenciesTable from "./CurrenciesTable";

export default function CurrenciesPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">מטבעות</h1>
      <CurrenciesTable />
    </div>
  );
}
