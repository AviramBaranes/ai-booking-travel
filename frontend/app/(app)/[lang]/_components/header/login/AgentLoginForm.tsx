"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { AlertCircle } from "lucide-react";
import { signIn } from "next-auth/react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { Loading } from "@/shared/components/Loading";

function schema(t: (key: string) => string) {
  return z.object({
    email: z.string().email(t("validation.invalidEmail")),
    password: z.string().min(1, t("validation.passwordRequired")),
  });
}
type FormData = z.infer<ReturnType<typeof schema>>;

const inputClass =
  "h-15 bg-background border-border-light rounded-xl px-6 text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:border-navy aria-invalid:bg-destructive/10";

interface Props {
  onSuccess: () => void;
}

export function AgentLoginForm({ onSuccess }: Props) {
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
      const result = await signIn("credentials", {
        redirect: false,
        email: data.email,
        password: data.password,
      });
      const res = result as { error?: string } | undefined;
      if (res?.error) throw new Error(res.error ?? "unknown_error");
      return result;
    },
    onSuccess: () => {
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
        type="submit"
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        disabled={isPending || !email?.trim() || !password?.trim()}
      >
        {isPending && <Loading />}
        {t("agent.submit")}
      </Button>
    </form>
  );
}
