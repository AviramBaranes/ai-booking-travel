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

interface CouponPopoverProps {
  checkboxLabel: string;
  inputLabel: string;
  saveButtonText: string;
  couponCode: string;
  setCouponCode: (code: string) => void;
}
export function CouponPopover({
  checkboxLabel,
  inputLabel,
  saveButtonText,
  couponCode,
  setCouponCode,
}: CouponPopoverProps) {
  const [hasCoupon, setHasCoupon] = useState(!!couponCode);
  const [isCouponSaved, setIsCouponSaved] = useState(false);

  return (
    <Popover open={hasCoupon && !isCouponSaved}>
      <PopoverTrigger asChild>
        <Field orientation="horizontal" className="w-auto shrink-0">
          <Checkbox
            checked={hasCoupon}
            onCheckedChange={(checked) => {
              setHasCoupon(!!checked);
              if (!checked) {
                setIsCouponSaved(false);
                setCouponCode("");
              }
            }}
            id="has-coupon"
            name="has-coupon"
            className="border-white w-3 h-3 rounded-xs bg-navy data-checked:bg-white data-checked:text-navy data-checked:border-white"
          />
          <FieldLabel htmlFor="has-coupon" className="text-white">
            {checkboxLabel}
          </FieldLabel>
        </Field>
      </PopoverTrigger>
      <PopoverContent className="py-2" align="start">
        <Field orientation="horizontal">
          <Input
            id="coupon"
            value={couponCode}
            onChange={(e) => {
              setCouponCode(e.target.value);
            }}
            placeholder={inputLabel}
            className="w-3/4 py-5 rounded-sm bg-background focus-visible:ring-0 focus-visible:border-transparent"
          />
          <Button
            variant="brand"
            onClick={() => setIsCouponSaved(true)}
            className="w-1/4 rounded-sm type-paragraph font-semibold py-5"
          >
            {saveButtonText}
          </Button>
        </Field>
      </PopoverContent>
    </Popover>
  );
}
