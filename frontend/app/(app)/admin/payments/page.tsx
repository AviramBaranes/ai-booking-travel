import { SupplierPaymentsShell } from "./_components/SupplierPaymentsShell";

export default function SupplierPaymentsPage() {
  return (
    <div className="mx-auto flex w-full max-w-7xl flex-col gap-6">
      <header className="flex flex-col gap-1">
        <h1 className="type-h4 text-navy">תשלומים לספקים</h1>
        <p className="type-paragraph text-text-secondary">
          בחר ספק כדי להציג את ההזמנות שטרם שולמו לו ולסמן אותן לתשלום.
        </p>
      </header>

      <SupplierPaymentsShell />
    </div>
  );
}
