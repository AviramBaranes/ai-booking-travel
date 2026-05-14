export function AccountCardLabelValue({
  label,
  value,
  valClassName,
}: {
  label: string;
  value: string;
  valClassName?: string;
}) {
  return (
    <div className="px-6 py-1 flex flex-col">
      <p className="text-xs text-muted">{label}</p>
      <p className={`text-sm text-navy ${valClassName}`}>{value}</p>
    </div>
  );
}