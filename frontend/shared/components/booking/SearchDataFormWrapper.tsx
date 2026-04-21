import { useDirection } from "@/shared/hooks/useDirection";
import { clsx } from "clsx";
import { X } from "lucide-react";

export function SearchDataFormWrapper({
  onClose,
  children,
}: {
  onClose: () => void;
  children: React.ReactNode;
}) {
  const dir = useDirection();

  return (
    <div className="relative bg-navy py-4 px-2 rounded-xl">
      <X
        className={clsx("absolute top-2 cursor-pointer text-muted", {
          "left-2": dir === "rtl",
          "right-2": dir === "ltr",
        })}
        onClick={onClose}
      />
      {children}
    </div>
  );
}
