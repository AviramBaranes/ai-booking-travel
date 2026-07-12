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
import useAuthStore, { UserRole } from "@/shared/auth/authStore";
import { login } from "@/shared/api/accounts-api";

function schema(t: (key: string) => string) {
  return z.object({
    email: z.string().email(t("validation.invalidEmail")),
    password: z.string().min(1, t("validation.passwordRequired")),
  });
}
type FormData = z.infer<ReturnType<typeof schema>>;

export const inputClass =
  "lg:h-15 h-12 bg-background border-border-light rounded-xl px-6 text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:border-navy aria-invalid:bg-destructive/10";

interface AgentLoginFormProps {
  onSuccess: () => void;
  onForgotPassword: () => void;
}

export function AgentLoginForm({
  onSuccess,
  onForgotPassword,
}: AgentLoginFormProps) {
  const t = useTranslations("Login");
  const tError = useTranslations("ApiErrors");
  const store = useAuthStore();

  const {
    register,
    watch,
    handleSubmit,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema(t)) });

  const { mutate, error, isPending } = useMutation({
    mutationFn: async (data: FormData) => login(data.email, data.password),
    onSuccess: (response) => {
      store.setSession(response.accessToken, response.accessTokenExpiresAt, {
        id: response.id,
        email: response.email,
        firstName: response.firstName,
        lastName: response.lastName,
        role: response.role as UserRole,
        phoneNumber: response.phoneNumber,
        officeId: response.officeId,
        isAdminAsAgent: false,
      });
      onSuccess();
    },
  });

  const hasError = !!error;
  const email = watch("email");
  const password = watch("password");

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

      <div>
        <Input
          type="password"
          placeholder={t("agent.password")}
          aria-invalid={hasError || !!errors.password}
          className={inputClass}
          {...register("password")}
        />
        <ErrorDisplay>{errors.password?.message}</ErrorDisplay>
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
        type="button"
        variant="ghost"
        className="w-1/3 underline mx-auto"
        onClick={onForgotPassword}
      >
        {t("agent.forgotPassword")}
      </Button>

      <Button
        type="submit"
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        loading={isPending}
        disabled={isPending || !email?.trim() || !password?.trim()}
      >
        {t("agent.submit")}
      </Button>
    </form>
  );
}
