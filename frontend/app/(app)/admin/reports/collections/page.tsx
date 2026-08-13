import { adminGate } from "@/lib/admin-gate";
import CollectionsReportTable from "./CollectionsReportTable";

export default async function CollectionsReportPage() {
  const gate = await adminGate("/admin/reports/collections");
  if (gate) return gate;
  return (
    <div className="space-y-4">
      <h1 className="text-xl font-bold text-navy">דוח גבייה</h1>
      <CollectionsReportTable />
    </div>
  );
}
