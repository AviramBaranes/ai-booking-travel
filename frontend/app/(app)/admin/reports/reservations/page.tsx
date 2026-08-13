import { adminGate } from "@/lib/admin-gate";
import ReservationsReportTable from "./ReservationsReportTable";

export default async function ReservationsReportPage() {
  const gate = await adminGate("/admin/reports/reservations");
  if (gate) return gate;
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">דוח הזמנות עסקי</h1>
      <ReservationsReportTable />
    </div>
  );
}
