import { BillingShell } from "./_components/BillingShell";

export default function BillingPage() {
  return (
    <div className="mx-auto flex w-full max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <h1 className="type-h4 text-navy">חיוב חשבונות פתוחים</h1>
        <p className="type-paragraph text-text-secondary">
          בחר רשת או סוכנות כדי להציג את ההזמנות הפתוחות שלה ולהפיק חיוב.
        </p>
      </header>

      <BillingShell />
    </div>
  );
}
