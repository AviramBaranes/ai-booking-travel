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
    <div className="flex flex-col items-center gap-4 text-center">
      {pillText && (
        <span className="rounded-full bg-brand-blue/10 min-w-50 py-3 type-label text-brand-blue">
          {pillText}
        </span>
      )}

      <h3 className="type-h5 font-extrabold lg:type-h3 text-navy">{title}</h3>

      {subtitle && <p className="type-h6 text-muted">{subtitle}</p>}
    </div>
  );
}
