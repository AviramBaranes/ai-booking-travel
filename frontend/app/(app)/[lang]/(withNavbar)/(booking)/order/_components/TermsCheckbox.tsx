import { Checkbox } from "@/components/ui/checkbox";
import { OrderFormValues } from "@/shared/components/booking/OrderForm/orderFormSchema";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useTranslations } from "next-intl";
import Link from "next/link";
import { Controller, useFormContext } from "react-hook-form";

interface TermsCheckboxProps {
  href: string;
}
export function TermsCheckbox({ href }: TermsCheckboxProps) {
  const t = useTranslations("booking.orderPage");
  const {
    control,
    formState: { errors },
  } = useFormContext<OrderFormValues>();

  return (
    <Controller
      name="termsAccepted"
      control={control}
      render={({ field }) => (
        <>
          <label className="flex items-center gap-2 cursor-pointer text-navy mx-auto max-sm:mb-4">
            <Checkbox
              checked={field.value}
              onCheckedChange={field.onChange}
              aria-invalid={!!errors?.termsAccepted}
              aria-errormessage={errors?.termsAccepted?.message}
              className="border-[#a9a8b3] data-checked:border-brand data-checked:bg-brand"
            />
            <span className="type-paragraph text-navy">
              {t("termsCheckbox")}{" "}
              <Link
                target="_blank"
                href={href}
                className="text-link underline type-label"
              >
                {t("termsLink")}
              </Link>
            </span>
          </label>
          <ErrorDisplay>{errors?.termsAccepted?.message}</ErrorDisplay>
        </>
      )}
    />
  );
}
