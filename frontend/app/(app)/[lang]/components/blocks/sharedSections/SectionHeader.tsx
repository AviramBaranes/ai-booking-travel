type SectionHeaderProps = {
  pillText?: string | null;
  title: string;
  subtitle?: string | null;
};

export function SectionHeader({
  pillText,
  title,
  subtitle,
}: SectionHeaderProps) {
  return (
    <div className="flex flex-col items-center gap-4">
      {pillText && (
        <span className="rounded-full bg-brand-blue/10 px-8 py-3 text-sm font-bold text-brand-blue">
          {pillText}
        </span>
      )}
      <h2 className="text-4xl font-black text-navy">{title}</h2>
      {subtitle && <p className="text-base text-muted">{subtitle}</p>}
    </div>
  );
}
