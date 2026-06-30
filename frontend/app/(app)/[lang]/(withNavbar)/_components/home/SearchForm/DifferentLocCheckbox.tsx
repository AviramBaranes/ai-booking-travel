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
        className="border-navy border-2 lg:border lg:border-white w-4 h-4 lg:w-3 lg:h-3 rounded-xs bg-transparent lg:bg-navy data-checked:bg-white data-checked:text-navy lg:data-checked:border-white"
      />
      <FieldLabel htmlFor="dropoff-different-loc" className="text-navy lg:text-white">
        {label}
      </FieldLabel>
    </Field>
  );
}
