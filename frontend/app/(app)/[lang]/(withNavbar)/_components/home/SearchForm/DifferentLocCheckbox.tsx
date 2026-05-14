import { Field, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";

interface SearchFormOptionsProps {
  label: string;
  isReturnDifferentLoc: boolean;
  setIsReturnDifferentLoc: (value: boolean) => void;
}

export function DifferentLocCheckbox({
  label,
  isReturnDifferentLoc,
  setIsReturnDifferentLoc,
}: SearchFormOptionsProps) {
  return (
    <Field orientation="horizontal" className="w-auto shrink-0">
      <Checkbox
        checked={isReturnDifferentLoc}
        onCheckedChange={(checked) => setIsReturnDifferentLoc(!!checked)}
        id="return-different-loc"
        name="return-different-loc"
        className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
      />
      <FieldLabel htmlFor="return-different-loc" className="text-white">
        {label}
      </FieldLabel>
    </Field>
  );
}
