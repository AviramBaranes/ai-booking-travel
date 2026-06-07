"use client";

import { useState } from "react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { AlertCircle, CheckCircle2 } from "lucide-react";
import { z } from "zod";
import Link from "next/link";
import { useParams, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { resetPassword } from "@/shared/api/accounts-api";

function schema(t: (key: string) => string) {
  return z
    .object({
      password: z.string().min(1, t("validation.passwordRequired")),
      confirmPassword: z.string().min(1, t("validation.passwordRequired")),
    })
    .refine((data) => data.password === data.confirmPassword, {
      message: t("validation.passwordsMustMatch"),
      path: ["confirmPassword"],
    });
}

type FormData = z.infer<ReturnType<typeof schema>>;

const inputClass =
  "h-15 bg-background border-border-light rounded-xl px-6 text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:border-navy aria-invalid:bg-destructive/10";

export function PasswordResetForm() {
  const t = useTranslations("PasswordReset");
  const params = useParams();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");
  const [isSuccess, setIsSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema(t)) });

  const { mutate, error, isPending } = useMutation({
    mutationFn: (data: FormData) => {
      if (!token) throw new Error("invalid_reset_token");
      return resetPassword({ token, newPassword: data.password });
    },
    onSuccess: () => {
      setIsSuccess(true);
    },
  });

  const onSubmit = (data: FormData) => {
    mutate(data);
  };

  if (isSuccess) {
    return (
      <div className="flex flex-col items-center gap-4 text-center">
        <CheckCircle2 className="size-12 text-primary" />
        <h2 className="text-2xl font-bold text-navy">
          {t("successTitle")}
        </h2>
        <Button variant="brand" className="w-full mt-4" asChild>
          <Link href={`/${params.lang}`}>
            {t("backToHome")}
          </Link>
        </Button>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-3 w-full">
      <div className="text-center mb-6">
        <h1 className="text-2xl font-bold text-navy mb-2">{t("title")}</h1>
        <p className="text-text-secondary">{t("subtitle")}</p>
      </div>

      <div>
        <Input
          type="password"
          placeholder={t("newPassword")}
          aria-invalid={!!errors.password}
          className={inputClass}
          {...register("password")}
        />
        <ErrorDisplay>{errors.password?.message}</ErrorDisplay>
      </div>

      <div>
        <Input
          type="password"
          placeholder={t("confirmPassword")}
          aria-invalid={!!errors.confirmPassword}
          className={inputClass}
          {...register("confirmPassword")}
        />
        <ErrorDisplay>{errors.confirmPassword?.message}</ErrorDisplay>
      </div>


      {(!token || !!error) && (
        <div role="alert" className="flex items-center gap-1">
          <AlertCircle className="size-3.5 text-destructive shrink-0" />
          <span className="type-paragraph text-destructive">
            {t("invalid_reset_token")}
          </span>
        </div>
      )}

      <Button
        type="submit"
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        loading={isPending}
        disabled={isPending || !token}
      >
        {t("submit")}
      </Button>
    </form>
  );
}
