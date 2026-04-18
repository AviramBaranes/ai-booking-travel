import { clsx } from "clsx";

export function Bone({
  className,
  variant = "light",
}: {
  className?: string;
  variant?: "light" | "dark";
}) {
  return (
    <div
      className={clsx(
        "animate-pulse rounded-md",
        {
          "bg-gray-200": variant === "light",
          "bg-white/50": variant === "dark",
        },
        className,
      )}
    />
  );
}
