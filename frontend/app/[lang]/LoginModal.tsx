"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { X } from "lucide-react";
import { getSession, signIn } from "next-auth/react";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/shared/components/Button";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";

function loginSchema(t: (key: string) => string) {
  return z.object({
    username: z.string().min(1, t("validation.usernameRequired")),
    password: z.string().min(1, t("validation.passwordRequired")),
  });
}

type LoginFormData = z.infer<ReturnType<typeof loginSchema>>;

export function LoginModal() {
  const t = useTranslations("Login");
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

  const loginMutation = useMutation({
    mutationFn: async (data: LoginFormData) => {
      const result = await signIn("credentials", {
        redirect: false,
        username: data.username,
        password: data.password,
      });
      if (!result?.ok) {
        throw new Error("Login failed");
      }
      console.log({ result });
      return result;
    },
    onSuccess: async () => {
      const session = await getSession();
      console.log("logging the session now:");
      console.log({ session });
      reset();
      if (session?.user?.role === "admin") {
        router.push("/admin");
      } else {
        setShowModal(false);
      }
    },
  });

  const onSubmit = (data: LoginFormData) => {
    loginMutation.mutate(data);
  };

  const closeModal = () => {
    reset();
    setShowModal(false);
  };

  return (
    <>
      <button className="cursor-pointer" onClick={() => setShowModal(true)}>
        {t("openModal")}
      </button>
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
              onSubmit={handleSubmit(onSubmit)}
              className="flex flex-col gap-3"
            >
              <div>
                <input
                  type="text"
                  placeholder={t("username")}
                  className="border p-2 rounded w-full"
                  {...register("username")}
                />
                <ErrorDisplay>{errors.username?.message}</ErrorDisplay>
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
              <Button
                type="submit"
                loading={loginMutation.isPending}
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
