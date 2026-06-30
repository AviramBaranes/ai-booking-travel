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
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslations } from "next-intl";
import { SearchFormCheckbox } from "./SearchFormCheckbox";

interface AgePopoverProps {
  saveButtonText: string;
  driverAge: number;
  setDriverAge: (age: number) => void;
}

export function AgePopover({
  saveButtonText,
  driverAge,
  setDriverAge,
}: AgePopoverProps) {
  const t = useTranslations("SearchForm");

  const [isAgeApproved, setIsAgeApproved] = useState(true);
  const [isChangedAge, setIsChangedAge] = useState(driverAge !== 30);
  const [isValid, setIsValid] = useState(true);

  return (
    <Popover open={!isAgeApproved}>
      <PopoverTrigger asChild>
        <SearchFormCheckbox
          label={t("ageRange", {
            ageRange:
              isChangedAge && driverAge >= 18 && driverAge <= 99
                ? driverAge
                : "30 - 65",
          })}
          value={isAgeApproved}
          setValue={(checked) => {
            if (!checked) {
              setIsAgeApproved(false);
            }
          }}
          id="age-above-30"
          name="age-above-30"
        />
      </PopoverTrigger>
      <PopoverContent className="py-2 w-auto min-w-max" align="end">
        <Field orientation="horizontal" className="flex items-start">
          <FieldLabel htmlFor="age" className="w-fit whitespace-nowrap py-2">
            {t("agePopoverLabel")}
          </FieldLabel>
          <div className="flex flex-col items-start">
            <Input
              id="age"
              min={18}
              max={99}
              value={driverAge}
              onChange={(e) => {
                setDriverAge(Number(e.target.value));
              }}
              type="number"
              aria-invalid={!isValid}
              className="w-20 py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
            />
            {!isValid && <ErrorDisplay>18-99</ErrorDisplay>}
          </div>
          <Button
            variant="brand"
            onClick={() => {
              if (driverAge < 18 || driverAge > 99) {
                setIsValid(false);
                return;
              }
              setIsValid(true);
              setIsChangedAge(true);
              setIsAgeApproved(true);
            }}
            className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
          >
            {saveButtonText}
          </Button>
        </Field>
      </PopoverContent>
    </Popover>
  );
}
