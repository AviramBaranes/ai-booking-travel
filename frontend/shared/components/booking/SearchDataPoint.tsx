import { LucideIcon } from "lucide-react";

interface SearchDataPointProps {
  icon: LucideIcon;
  label: string;
  location: string;
  date: Date;
  time: string;
}

function formatDateTime(date: Date, time: string) {
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = date.getFullYear();
  return `${day}.${month}.${year} | ${time}`;
}

export function SearchDataPoint({
  icon: Icon,
  label,
  location,
  date,
  time,
}: SearchDataPointProps) {
  return (
    <div className="flex w-75 flex-col items-start gap-2">
      <span className="inline-flex items-center gap-2 rounded-full border border-white px-3 py-0.5">
        <Icon className="size-3 text-white" />
        <span className="type-paragraph text-white">{label}</span>
      </span>
      <p className="type-h6 w-full text-white">{location}</p>
      <p className="type-paragraph w-full text-white">
        {formatDateTime(date, time)}
      </p>
    </div>
  );
}
