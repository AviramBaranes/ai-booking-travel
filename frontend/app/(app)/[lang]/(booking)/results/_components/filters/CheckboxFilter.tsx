import type { LucideIcon } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";

export interface CheckboxFilterProps {
  title: string;
  icon: LucideIcon;
  options: { label: string; value: string }[];
  selectedValues: Set<string>;
  onChange: (value: string) => void;
  showDivider?: boolean;
}

export function CheckboxFilter({
  title,
  icon: Icon,
  options,
  selectedValues,
  onChange,
  showDivider = true,
}: CheckboxFilterProps) {
  return (
    <section
      className={showDivider ? "mt-5 border-t border-cars-border pt-5" : ""}
    >
      <h6 className="type-h6 mb-8 text-navy flex items-center gap-2">
        <Icon size={20} className="text-brand" />
        {title}
      </h6>
      <div className="flex flex-col gap-2">
        {options.map((option) => (
          <label
            key={option.value}
            className="flex items-center gap-2 cursor-pointer"
          >
            <Checkbox
              checked={selectedValues.has(option.value)}
              onCheckedChange={() => onChange(option.value)}
              className="border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand"
            />
            <span
              className={
                "type-paragraph text-navy " +
                (selectedValues.has(option.value) ? "font-bold" : "")
              }
            >
              {option.label}
            </span>
          </label>
        ))}
      </div>
    </section>
  );
}
