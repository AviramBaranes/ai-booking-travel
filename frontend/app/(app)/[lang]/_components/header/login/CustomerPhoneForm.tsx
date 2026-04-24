"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { AlertCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon } from "@/components/ui/input-group";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { Loading } from "@/shared/components/Loading";
import { sendOTP } from "@/shared/api/accounts-api";

const ISRAELI_MOBILE_RE = /^0[5][0-9][\s-]?\d{3}[\s-]?\d{4}$/;

function schema(t: (key: string) => string) {
  return z.object({
    phone: z
      .string()
      .min(1, t("validation.phoneRequired"))
      .regex(ISRAELI_MOBILE_RE, t("validation.invalidPhone")),
  });
}
type FormData = z.infer<ReturnType<typeof schema>>;

interface Props {
  onSubmit: (phone: string) => void;
}

export function CustomerPhoneForm({ onSubmit }: Props) {
  const t = useTranslations("Login");
  const tError = useTranslations("ApiErrors");

  const {
    register,
    watch,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema(t)) });

  const { mutate, error, isPending } = useMutation({
    mutationFn: async (data: FormData) => {
      const digits = data.phone.replace(/[\s-]/g, "");
      await sendOTP({ phoneNumber:digits });
      return data.phone;
    },
    onSuccess: (phone) => {
      onSubmit(phone);
    },
  });

  const phone = watch("phone");

  return (
    <form
      onSubmit={handleSubmit((d) => mutate(d))}
      className="flex flex-col gap-3 w-full"
    >
      <div>
        <InputGroup
          dir="ltr"
          className="h-15 rounded-xl bg-background border-border-light"
        >
          <InputGroupAddon align="inline-start" className="gap-3 ps-6">
            <div className="h-7.5 w-px bg-border-light/60" />
          </InputGroupAddon>
          <Input
            type="tel"
            placeholder={t("customer.phonePlaceholder")}
            aria-invalid={!!errors.phone}
            className="h-auto border-0 bg-transparent text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:ring-0 pe-6"
            {...register("phone")}
          />
        </InputGroup>
        <ErrorDisplay>{errors.phone?.message}</ErrorDisplay>
      </div>

      {error && (
        <div role="alert" className="flex items-center gap-1">
          <AlertCircle className="size-3.5 text-destructive shrink-0" />
          <span className="type-paragraph text-destructive">
            {tError(error.message)}
          </span>
        </div>
      )}

      <Button
        type="submit"
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        loading={isPending}
        disabled={isPending || !phone?.trim()}
      >
        {isPending && <Loading />}
        {t("customer.submit")}
      </Button>
    </form>
  );
}
