import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";
import type { ComponentProps } from "react";

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

  const handleNavigate: ComponentProps<typeof Link>["onNavigate"] = () => {
    window.scrollTo({ top: 0, left: 0, behavior: "auto" });
  };

  const params = new URLSearchParams(searchParams.toString());
  params.set("cid", String(carIndex));

  return (
    <Link
      href={`/${lang}/plans?${params.toString()}`}
      className={className}
      scroll
      onNavigate={handleNavigate}
    >
      {children}
    </Link>
  );
}
