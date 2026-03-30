import CouponsTable from "./CouponsTable";

export default function CouponsPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">קופונים</h1>
      <CouponsTable />
    </div>
  );
}
