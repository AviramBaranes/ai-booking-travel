"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { X, User } from "lucide-react";
import { getSession, signIn } from "next-auth/react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { Button } from "@/components/ui/button";

function loginSchema(t: (key: string) => string) {
  return z.object({
    email: z.string().email(t("validation.invalidEmail")),
    password: z.string().min(1, t("validation.passwordRequired")),
  });
}

type LoginFormData = z.infer<ReturnType<typeof loginSchema>>;

export function LoginModal() {
  const t = useTranslations("Login");
  const tError = useTranslations("ApiErrors");
  const router = useRouter();
  const [showModal, setShowModal] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema(t)),
  });

  const { mutate, error, isPending } = useMutation({
    mutationFn: async (data: LoginFormData) => {
      const result = await signIn("credentials", {
        redirect: false,
        email: data.email,
        password: data.password,
      });

      if (result?.error) {
        throw new Error(result?.error ?? "unknown_error");
      }

      return result;
    },
    onSuccess: async () => {
      const session = await getSession();
      reset();
      if (session?.user?.role === "admin") {
        router.push("/admin");
      } else {
        setShowModal(false);
      }
    },
  });

  const closeModal = () => {
    reset();
    setShowModal(false);
  };

  return (
    <>
      <Button
        size="outline"
        variant="outline"
        onClick={() => setShowModal(true)}
      >
        <User className="size-5" />
        {t("openModal")}
      </Button>
      {showModal && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 cursor-pointer"
          onClick={closeModal}
        >
          <div
            className="relative bg-white p-6 rounded-lg shadow-lg w-80 cursor-default"
            onClick={(e) => e.stopPropagation()}
          >
            <button
              onClick={closeModal}
              className="absolute top-3 inset-e-3 text-gray-400 hover:text-gray-600 cursor-pointer"
            >
              <X size={20} />
            </button>
            <h2 className="text-xl mb-4">{t("title")}</h2>
            <form
              onSubmit={handleSubmit((d) => mutate(d))}
              className="flex flex-col gap-3"
            >
              <div>
                <input
                  type="email"
                  placeholder={t("email")}
                  className="border p-2 rounded w-full"
                  {...register("email")}
                />
                <ErrorDisplay>{errors.email?.message}</ErrorDisplay>
              </div>
              <div>
                <input
                  type="password"
                  placeholder={t("password")}
                  className="border p-2 rounded w-full"
                  {...register("password")}
                />
                <ErrorDisplay>{errors.password?.message}</ErrorDisplay>
              </div>
              <ErrorDisplay>{error && tError(error.message)}</ErrorDisplay>
              <Button
                type="submit"
                // loading={isPending}
                className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {t("submit")}
              </Button>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
