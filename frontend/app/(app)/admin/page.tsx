import ReservationsReportTable from "./reports/reservations/ReservationsReportTable";

function toDateInputValue(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

export default function AdminHome() {
  const today = new Date();
  const pickupDateFrom = toDateInputValue(today);
  const pickupDateTo = toDateInputValue(
    new Date(today.getFullYear(), today.getMonth(), today.getDate() + 4),
  );

  return (
    <div className="space-y-4">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold text-gray-700">
          הזמנות דורשות התייחסות
        </h1>
        <p className="text-sm text-gray-500">
          הזמנות שלא כורטסו והאיסוף שלהם הוא בעוד 4 או פחות ימים.
        </p>
      </div>
      <ReservationsReportTable
        showFilters={false}
        fixedFilters={{
          status: "booked",
          pickupFrom: pickupDateFrom,
          pickupTo: pickupDateTo,
        }}
      />
    </div>
  );
}
