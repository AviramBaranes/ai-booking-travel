"use client";
import { Button } from "@/components/ui/button";
import { updateMe } from "@/shared/api/accounts-api";
import useAuthStore from "@/shared/auth/authStore";
import { CustomerForm } from "@/shared/components/booking/OrderForm/CustomerForm";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { Loading } from "@/shared/components/Loading";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { zodResolver } from "@hookform/resolvers/zod/dist/zod.js";
import { useMutation } from "@tanstack/react-query";
import { useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import zod from "zod";

function useUpdateCustomerSchema() {
  const t = useTranslations("MyAccount.profile.validation");
  return zod.object({
    customerFirstName: zod.string().min(1, t("required")),
    customerLastName: zod.string().min(1, t("required")),
    customerPhone: zod.string().min(1, t("required")),
    customerEmail: zod.email(t("invalidEmail")).min(1, t("required")),
  });
}
type UpdateCustomer = {
  customerFirstName: string;
  customerLastName: string;
  customerPhone: string;
  customerEmail: string;
};

export function ProfilePageContent() {
  const t = useTranslations("MyAccount.profile");
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const status = useAuthStore((state) => state.status);
  const updateCustomerSchema = useUpdateCustomerSchema();

  const isLoading = status === "loading" || status === "idle";
  const isAuthorized = !isLoading && user?.role === "customer";

  useEffect(() => {
    if (isLoading) {
      return;
    }

    if (!isAuthorized) {
      router.replace("/he/");
    }
  }, [isAuthorized, router, isLoading]);

  const formMethods = useForm<UpdateCustomer>({
    resolver: zodResolver(updateCustomerSchema),
    defaultValues: {
      customerFirstName: user?.firstName ?? "",
      customerLastName: user?.lastName ?? "",
      customerPhone: user?.phoneNumber ?? "",
      customerEmail: user?.email ?? "",
    },
  });

  const { handleSubmit, reset } = formMethods;

  useEffect(() => {
    if (!user) return;

    // Keep already edited inputs untouched while filling in the rest.
    reset(
      {
        customerFirstName: user.firstName ?? "",
        customerLastName: user.lastName ?? "",
        customerPhone: user.phoneNumber ?? "",
        customerEmail: user.email ?? "",
      },
      { keepDirtyValues: true },
    );
  }, [user, reset]);

  const { mutate, isPending, isSuccess, error } = useMutation({
    mutationFn: async (data: UpdateCustomer) =>
      updateMe({
        firstName: data.customerFirstName,
        lastName: data.customerLastName,
        phoneNumber: data.customerPhone,
        email: data.customerEmail,
      }),
  });

  const translatedError = useTranslatedError(error);

  if (isLoading || !isAuthorized) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loading />
      </div>
    );
  }

  return (
    <main className="lg:w-2/3 lg:mx-auto">
      <h5 className="type-h5 text-navy my-6">{t("title")}</h5>
      <form onSubmit={handleSubmit((data) => mutate(data))}>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <FormProvider {...formMethods}>
            <CustomerForm />
          </FormProvider>
        </div>
        <Button
          variant="brand"
          loading={isPending}
          type="submit"
          className="my-5 font-bold px-10 py-5 rounded-md"
        >
          {t("saveButton")}
        </Button>
        {isSuccess && <SuccessBadge>{t("successMessage")}</SuccessBadge>}
        <ErrorDisplay>{translatedError}</ErrorDisplay>
      </form>
    </main>
  );
}
