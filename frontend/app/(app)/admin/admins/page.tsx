import AdminsTable from "./AdminsTable";

export default function AdminsPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">מנהלים</h1>
      <AdminsTable />
    </div>
  );
}
