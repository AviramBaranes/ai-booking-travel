import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { ClassAttributes, InputHTMLAttributes } from "react";

export function OrderFormInput({
  className,
  ...props
}: ClassAttributes<HTMLInputElement> & InputHTMLAttributes<HTMLInputElement>) {
  const readonlyClass = props.readOnly ? "bg-muted/10 cursor-not-allowed" : "bg-white";

  return (
    <Input
      className={cn(
        readonlyClass,
        "border border-cars-border h-12 rounded-lg px-4 type-paragraph text-text-secondary w-full",
        className,
      )}
      {...props}
    />
  );
}
