import { Field, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";

interface SearchFormCheckboxProps
  extends Omit<React.ComponentPropsWithRef<typeof Field>, "children"> {
  label: string;
  value: boolean;
  setValue: (value: boolean) => void;
  id?: string;
  name?: string;
}

export function SearchFormCheckbox({
  label,
  value,
  setValue,
  id,
  name,
  className,
  ref,
  ...props
}: SearchFormCheckboxProps) {
  return (
    <Field
      ref={ref}
      orientation="horizontal"
      className={cn("w-auto shrink-0", className)}
      {...props}
    >
      <Checkbox
        checked={value}
        onCheckedChange={(checked) => setValue(!!checked)}
        id={id}
        name={name}
        className="h-4 w-4 rounded-xs border-2 data-checked:border-navy bg-white data-checked:bg-white data-checked:text-navy lg:h-3 lg:w-3 lg:border lg:border-white lg:bg-navy lg:data-checked:border-white"
      />

      <FieldLabel htmlFor={id} className="text-navy lg:text-white">
        {label}
      </FieldLabel>
    </Field>
  );
}