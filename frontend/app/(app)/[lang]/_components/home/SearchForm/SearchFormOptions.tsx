import { useState } from "react";
import { useTranslations } from "next-intl";
import { Field, FieldLabel } from "@/components/ui/field";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface SearchFormOptionsProps {
  isReturnDifferentLoc: boolean;
  setIsReturnDifferentLoc: (value: boolean) => void;
}

export function SearchFormOptions({
  isReturnDifferentLoc,
  setIsReturnDifferentLoc,
}: SearchFormOptionsProps) {
  const t = useTranslations("SearchFormOptions");
  const [isAgeNormal, setIsAgeNormal] = useState(true);
  const [isChangedAge, setIsChangedAge] = useState(false);
  const [coupon, setCoupon] = useState("");
  const [hasCoupon, setHasCoupon] = useState(false);
  const [isCouponSaved, setIsCouponSaved] = useState(false);
  const [driverAge, setDriverAge] = useState(30);

  return (
    <div className="bg-navy w-fit h-9 rounded-t-xl flex items-center text-white type-h6 px-6 gap-5">
      <Field orientation="horizontal" className="w-auto shrink-0">
        <Checkbox
          checked={isReturnDifferentLoc}
          onCheckedChange={(checked) => setIsReturnDifferentLoc(!!checked)}
          id="return-different-loc"
          name="return-different-loc"
          className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
        />
        <FieldLabel htmlFor="return-different-loc" className="text-white">
          {t("returnDifferentLoc")}
        </FieldLabel>
      </Field>
      <div className="h-4 w-px bg-white/40 shrink-0" />
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
              {t("ageRange")}
            </FieldLabel>
          </Field>
        </PopoverTrigger>
        <PopoverContent className="py-2 w-auto min-w-max">
          <Field orientation="horizontal">
            <FieldLabel htmlFor="age" className="w-fit whitespace-nowrap">
              {t("agePopoverLabel")}
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
              {t("save")}
            </Button>
          </Field>
        </PopoverContent>
      </Popover>
      <div className="h-4 w-px bg-white/40 shrink-0" />
      <Popover open={hasCoupon && !isCouponSaved}>
        <PopoverTrigger asChild>
          <Field orientation="horizontal" className="w-auto shrink-0">
            <Checkbox
              checked={hasCoupon}
              onCheckedChange={(checked) => {
                setHasCoupon(!!checked);
                if (!checked) {
                  setIsCouponSaved(false);
                  setCoupon("");
                }
              }}
              id="has-coupon"
              name="has-coupon"
              className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
            />
            <FieldLabel htmlFor="has-coupon" className="text-white">
              {t("hasCoupon")}
            </FieldLabel>
          </Field>
        </PopoverTrigger>
        <PopoverContent className="py-2">
          <Field orientation="horizontal">
            <Input
              id="coupon"
              value={coupon}
              onChange={(e) => {
                setCoupon(e.target.value);
              }}
              placeholder={t("couponPlaceholder")}
              className="w-3/4 py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
            />
            <Button
              variant="brand"
              onClick={() => setIsCouponSaved(true)}
              className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
            >
              {t("save")}
            </Button>
          </Field>
        </PopoverContent>
      </Popover>
    </div>
  );
}
