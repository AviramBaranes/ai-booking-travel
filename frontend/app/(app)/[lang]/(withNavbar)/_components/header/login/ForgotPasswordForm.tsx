"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { AlertCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { sendPasswordReset } from "@/shared/api/accounts-api";

function schema(t: (key: string) => string) {
  return z.object({
    email: z.string().email(t("validation.invalidEmail")),
  });
}
type FormData = z.infer<ReturnType<typeof schema>>;

const inputClass =
  "h-15 bg-background border-border-light rounded-xl px-6 text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:border-navy aria-invalid:bg-destructive/10";

interface ForgotPasswordFormProps {
  onSuccess: () => void;
  onBackToLogin: () => void;
}

export function ForgotPasswordForm({
  onSuccess,
  onBackToLogin,
}: ForgotPasswordFormProps) {
  const t = useTranslations("Login");
  const tError = useTranslations("ApiErrors");

  const {
    register,
    watch,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema(t)) });

  const { mutate, error, isPending } = useMutation({
    mutationFn: async (data: FormData) => sendPasswordReset(data),
    onSuccess: () => {
      onSuccess();
    },
  });

  const hasError = !!error;
  const email = watch("email");

  return (
    <form
      onSubmit={handleSubmit((d) => mutate(d))}
      className="flex flex-col gap-3 w-full"
    >
      <div>
        <Input
          type="email"
          placeholder={t("agent.email")}
          aria-invalid={!!errors.email}
          className={inputClass}
          {...register("email")}
        />
        <ErrorDisplay>{errors.email?.message}</ErrorDisplay>
      </div>

      {hasError && (
        <div role="alert" className="flex items-center gap-1">
          <AlertCircle className="size-3.5 text-destructive shrink-0" />
          <span className="type-paragraph text-destructive">
            {tError(error.message)}
          </span>
        </div>
      )}

      <Button
        variant="ghost"
        type="button"
        className="w-1/3 underline mx-auto"
        onClick={onBackToLogin}
      >
        {t("agent.backToLogin")}
      </Button>

      <Button
        type="submit"
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        loading={isPending}
        disabled={isPending || !email?.trim()}
      >
        {t("agent.sendPasswordReset")}
      </Button>
    </form>
  );
}
