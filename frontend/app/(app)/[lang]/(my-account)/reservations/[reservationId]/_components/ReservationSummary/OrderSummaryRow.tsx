export function OrderSummaryRow({
  label,
  value,
  valClassName,
}: {
  label: string;
  value: string;
  valClassName?: string;
}) {
  return (
    <div className="flex justify-between">
      <span className="type-paragraph text-text-secondary">{label}</span>
      <span className={`type-paragraph ${valClassName ?? "text-value"}`}>
        {value}
      </span>
    </div>
  );
}
