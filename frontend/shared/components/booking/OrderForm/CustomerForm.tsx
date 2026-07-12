import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { OrderFormValues } from "./orderFormSchema";
import { useTranslations } from "next-intl";
import { useFormContext } from "react-hook-form";
import { OrderFormInput } from "./OrderFormInput";
import { useDirection } from "@/shared/hooks/useDirection";
import useAuthStore from "@/shared/auth/authStore";

export function CustomerForm() {
  const t = useTranslations("booking.orderPage");
  const status = useAuthStore((s) => s.status);
  const readonly = status === "authenticated";
  const dir = useDirection();
  const {
    register,
    formState: { errors },
  } = useFormContext<OrderFormValues>();

  return (
    <>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("customerFirstName")}
          aria-invalid={!!errors.customerFirstName}
          readOnly={readonly}
          {...register("customerFirstName")}
        />
        <ErrorDisplay>{errors.customerFirstName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("customerLastName")}
          aria-invalid={!!errors.customerLastName}
          readOnly={readonly}
          {...register("customerLastName")}
        />
        <ErrorDisplay>{errors.customerLastName?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("customerEmail")}
          type="email"
          aria-invalid={!!errors.customerEmail}
          readOnly={readonly}
          {...register("customerEmail")}
        />
        <ErrorDisplay>{errors.customerEmail?.message}</ErrorDisplay>
      </div>
      <div className="flex-1">
        <OrderFormInput
          placeholder={t("customerPhone")}
          type="tel"
          pattern="^05\d{8}$"
          dir={dir}
          aria-invalid={!!errors.customerPhone}
          readOnly={readonly}
          {...register("customerPhone")}
        />
        <ErrorDisplay>{errors.customerPhone?.message}</ErrorDisplay>
      </div>
    </>
  );
}
