import { adminGate } from "@/lib/admin-gate";
import { MonthlyReportsBrowser } from "./_components/MonthlyReportsBrowser";

export default async function MonthlyReportsPage() {
  const gate = await adminGate("/admin/monthly-reports");
  if (gate) return gate;
  return (
    <div className="mx-auto flex w-full max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <h1 className="type-h4 text-navy">דוחות חודשיים</h1>
        <p className="type-paragraph text-text-secondary">
          תיקייה לכל גורם מחייב, ובתוכה הדוחות החודשיים שנשלחו אליו.
        </p>
      </header>

      <MonthlyReportsBrowser />
    </div>
  );
}
