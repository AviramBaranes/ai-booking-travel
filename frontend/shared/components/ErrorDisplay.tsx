import clsx from "clsx";
import { PropsWithChildren } from "react";

export function ErrorDisplay({
  children,
  className,
}: PropsWithChildren<{ className?: string }>) {
  return (
    <label className={clsx("label", { hidden: !children }, className)}>
      <span className="label-text text-red-500 text-sm">{children}</span>
    </label>
  );
}
