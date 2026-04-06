import clsx from "clsx";
import { PropsWithChildren } from "react";

export function ErrorDisplay({
  children,
  className,
}: PropsWithChildren<{ className?: string }>) {
  if (!children) return null;
  return (
    <p
      role="alert"
      className={clsx("text-destructive text-sm mt-1", className)}
    >
      {children}
    </p>
  );
}
