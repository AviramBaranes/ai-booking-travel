import CustomersTable from "./CustomersTable";

export default function CustomersPage() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-700">לקוחות פרטיים</h1>
      <CustomersTable />
    </div>
  );
}