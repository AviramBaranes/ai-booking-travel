import { adminGate } from "@/lib/admin-gate";
import ProfitReportTable from "./ProfitReportTable";

export default async function ProfitabilityReportPage() {
  const gate = await adminGate("/admin/reports/profitability");
  if (gate) return gate;
  return (
    <div className="space-y-4">
      <h1 className="text-xl font-bold text-navy">דוח רווחיות</h1>
      <ProfitReportTable />
    </div>
  );
}
