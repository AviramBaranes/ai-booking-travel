import LocationsTable from "./LocationsTable";

export default function LocationsPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">מיקומים</h1>
      <LocationsTable />
    </div>
  );
}
