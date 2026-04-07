import { useState } from "react";
import { Field, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface AgePopoverProps {
  checkboxLabel: string;
  inputLabel: string;
  saveButtonText: string;
}

export function AgePopover({
  checkboxLabel,
  inputLabel,
  saveButtonText,
}: AgePopoverProps) {
  const [isAgeNormal, setIsAgeNormal] = useState(true);
  const [isChangedAge, setIsChangedAge] = useState(false);
  const [driverAge, setDriverAge] = useState(30);

  return (
    <Popover open={!isAgeNormal && !isChangedAge}>
      <PopoverTrigger asChild>
        <Field orientation="horizontal" className="w-auto shrink-0">
          <Checkbox
            checked={isAgeNormal}
            onCheckedChange={(checked) => {
              if (!checked) {
                setIsChangedAge(false);
                setDriverAge(30);
              }
              setIsAgeNormal(!!checked);
            }}
            id="age-above-30"
            name="age-above-30"
            className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
          />
          <FieldLabel htmlFor="age-above-30" className="text-white">
            {checkboxLabel}
          </FieldLabel>
        </Field>
      </PopoverTrigger>
      <PopoverContent className="py-2 w-auto min-w-max">
        <Field orientation="horizontal">
          <FieldLabel htmlFor="age" className="w-fit whitespace-nowrap">
            {inputLabel}
          </FieldLabel>
          <Input
            id="age"
            value={driverAge}
            onChange={(e) => {
              setDriverAge(Number(e.target.value));
            }}
            type="number"
            className="w-20 py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
          />
          <Button
            variant="brand"
            onClick={() => setIsChangedAge(true)}
            className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
          >
            {saveButtonText}
          </Button>
        </Field>
      </PopoverContent>
    </Popover>
  );
}
