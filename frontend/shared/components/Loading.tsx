import { LoaderCircle } from "lucide-react";

export function Loading({ className }: { className?: string }) {
  return <LoaderCircle className={`animate-spin ${className ?? ""}`} />;
}
