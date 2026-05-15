import { Field, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";

interface SearchFormOptionsProps {
  label: string;
  isDropoffDifferentLoc: boolean;
  setIsDropoffDifferentLoc: (value: boolean) => void;
}

export function DifferentLocCheckbox({
  label,
  isDropoffDifferentLoc,
  setIsDropoffDifferentLoc,
}: SearchFormOptionsProps) {
  return (
    <Field orientation="horizontal" className="w-auto shrink-0">
      <Checkbox
        checked={isDropoffDifferentLoc}
        onCheckedChange={(checked) => setIsDropoffDifferentLoc(!!checked)}
        id="dropoff-different-loc"
        name="dropoff-different-loc"
        className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
      />
      <FieldLabel htmlFor="dropoff-different-loc" className="text-white">
        {label}
      </FieldLabel>
    </Field>
  );
}
