import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";

export function ContinueToPlansLink({
  carIndex,
  children,
  className,
}: {
  carIndex: number;
  children: React.ReactNode;
  className?: string;
}) {
  const { lang } = useParams();
  const searchParams = useSearchParams();

  const params = new URLSearchParams(searchParams.toString());
  params.set("cid", String(carIndex));

  return (
    <Link
      href={`/${lang}/plans?${params.toString()}`}
      className={className}
      scroll
    >
      {children}
    </Link>
  );
}
