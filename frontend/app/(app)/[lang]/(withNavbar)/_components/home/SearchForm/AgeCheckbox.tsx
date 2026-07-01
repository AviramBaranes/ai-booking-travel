import { useState } from "react";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslations } from "next-intl";
import { SearchFormCheckbox } from "./SearchFormCheckbox";

interface AgeCheckboxProps {
  saveButtonText: string;
  driverAge: number;
  setDriverAge: (age: number) => void;
}

export function AgeCheckbox({
  saveButtonText,
  driverAge,
  setDriverAge,
}: AgeCheckboxProps) {
  const t = useTranslations("SearchForm");

  const [isAgeApproved, setIsAgeApproved] = useState(true);
  const [isChangedAge, setIsChangedAge] = useState(driverAge !== 30);
  const [isValid, setIsValid] = useState(true);

  return (
    <div className="flex flex-col items-start w-9/10 mx-auto gap-2 my-2">
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
      {!isAgeApproved && (
        <Field orientation="horizontal" className="flex items-start justify-between">
          <FieldLabel htmlFor="age" className="w-fit whitespace-nowrap py-2 flex-none! text-navy">
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
              className="py-2.5 px-3.5 w-20 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
            />
            {!isValid && <ErrorDisplay>18-99</ErrorDisplay>}
          </div>
          <Button
            type="button"
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
            className="rounded-sm type-paragraph font-semibold py-2.5 px-3.5"
          >
            {saveButtonText}
          </Button>
        </Field>
      )}
    </div>
  );
}
