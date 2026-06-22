import { AddLocationAliasButton } from "./_components/AddLocationAliasButton";
import LocationsTable from "./_components/LocationsTable";

export default function LocationsPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">מיקומים</h1>
      <AddLocationAliasButton />
      <LocationsTable />
    </div>
  );
}
