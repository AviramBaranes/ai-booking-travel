import { useState } from "react";
import { Field } from "@/components/ui/field";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { SearchFormCheckbox } from "./SearchFormCheckbox";

interface CouponCheckboxProps {
  checkboxLabel: string;
  inputLabel: string;
  saveButtonText: string;
  couponCode: string;
  setCouponCode: (code: string) => void;
}
export function CouponCheckbox({
  checkboxLabel,
  inputLabel,
  saveButtonText,
  couponCode,
  setCouponCode,
}: CouponCheckboxProps) {
  const [hasCoupon, setHasCoupon] = useState(!!couponCode);
  const [isCouponSaved, setIsCouponSaved] = useState(false);

  return (
    <div className="flex flex-col items-start w-9/10 mx-auto gap-2 my-2">
      <SearchFormCheckbox
        label={checkboxLabel}
        value={hasCoupon}
        setValue={(checked) => {
          setHasCoupon(!!checked);
          if (!checked) {
            setIsCouponSaved(false);
            setCouponCode("");
          }
        }}
        id="has-coupon"
        name="has-coupon"
      />
      {hasCoupon && !isCouponSaved && (
        <Field orientation="horizontal">
          <Input
            id="coupon"
            value={couponCode}
            onChange={(e) => {
              setCouponCode(e.target.value);
            }}
            placeholder={inputLabel}
            className="py-2.5 px-3.5  rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
          />
          <Button
            variant="brand"
            onClick={() => setIsCouponSaved(true)}
            className="rounded-sm type-paragraph font-semibold py-2.5 px-3.5"
          >
            {saveButtonText}
          </Button>
        </Field>
      )}
    </div>
  );
}
