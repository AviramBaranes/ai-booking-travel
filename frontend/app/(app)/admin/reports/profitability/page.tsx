import ProfitReportTable from "./ProfitReportTable";

export default function ProfitabilityReportPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-xl font-bold text-navy">דוח רווחיות</h1>
      <ProfitReportTable />
    </div>
  );
}
