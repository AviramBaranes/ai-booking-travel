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

export function SearchFormOptions() {
  const [isReturnDifferentLoc, setIsReturnDifferentLoc] = useState(false);
  const [isAgeNormal, setIsAgeNormal] = useState(true);
  const [isChangedAge, setIsChangedAge] = useState(false);
  const [coupon, setCoupon] = useState("");
  const [hasCoupon, setHasCoupon] = useState(false);
  const [isCouponSaved, setIsCouponSaved] = useState(false);
  const [driverAge, setDriverAge] = useState(30);

  return (
    <div className="bg-navy w-[38%] h-9 rounded-t-xl flex items-center justify-between text-white type-h6 px-6">
      <Field
        orientation="horizontal"
        className="w-auto shrink-0 border-l pl-5 border-white/40"
      >
        <Checkbox
          checked={isReturnDifferentLoc}
          onCheckedChange={(checked) => setIsReturnDifferentLoc(!!checked)}
          id="return-different-loc"
          name="return-different-loc"
          className="border-white  w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
        />
        <FieldLabel htmlFor="return-different-loc" className="text-white">
          החזרת הרכב במקום אחר?
        </FieldLabel>
      </Field>
      <Popover open={!isAgeNormal && !isChangedAge}>
        <PopoverTrigger asChild>
          <Field
            orientation="horizontal"
            className="w-auto shrink-0 border-l pl-5 border-white/40"
          >
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
              גיל הנהג/ת: 30 - 65
            </FieldLabel>
          </Field>
        </PopoverTrigger>
        <PopoverContent className="py-2">
          <Field orientation="horizontal">
            <FieldLabel htmlFor="age">מה גיל הנהג/ת</FieldLabel>
            <Input
              id="age"
              value={driverAge}
              onChange={(e) => {
                setDriverAge(Number(e.target.value));
              }}
              type="number"
              className="w-1/4  py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
            />
            <Button
              variant="brand"
              onClick={() => setIsChangedAge(true)}
              className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
            >
              שמירה
            </Button>
          </Field>
        </PopoverContent>
      </Popover>
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
              יש לי שובר
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
              placeholder="הזינו שובר"
              className="w-3/4  py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
            />
            <Button
              variant="brand"
              onClick={() => setIsCouponSaved(true)}
              className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
            >
              שמירה
            </Button>
          </Field>
        </PopoverContent>
      </Popover>
    </div>
  );
}
