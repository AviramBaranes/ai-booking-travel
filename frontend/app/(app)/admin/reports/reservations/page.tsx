import ReservationsReportTable from "./ReservationsReportTable";

export default function ReservationsReportPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">דוח הזמנות עסקי</h1>
      <ReservationsReportTable />
    </div>
  );
}
