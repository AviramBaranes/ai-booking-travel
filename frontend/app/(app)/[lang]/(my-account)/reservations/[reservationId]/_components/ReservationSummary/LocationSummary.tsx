import { useParams } from "next/navigation";

interface LocationDateTimeSummaryProps {
  title: string;
  locationName: string;
  date: string;
  time: string;
  linkText: string;
}
export function LocationDateTimeSummary({
  title,
  locationName,
  date,
  time,
  linkText,
}: LocationDateTimeSummaryProps) {
  const { lang } = useParams();
  return (
    <div className="flex flex-col gap-2 mt-2">
      <p className="type-label text-navy">{title}</p>
      <p className="type-paragraph text-text-secondary">{locationName}</p>
      <p className="type-paragraph text-text-secondary">
        {new Date(date).toLocaleDateString(lang)} | {time}
      </p>
      <p className="text-brand underline type-label print:hidden">{linkText}</p>
    </div>
  );
}
