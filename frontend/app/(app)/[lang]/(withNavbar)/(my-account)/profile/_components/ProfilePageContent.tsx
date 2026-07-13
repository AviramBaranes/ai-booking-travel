"use client";
import { Button } from "@/components/ui/button";
import { updateMe } from "@/shared/api/accounts-api";
import useAuthStore from "@/shared/auth/authStore";
import { CustomerForm } from "@/shared/components/booking/OrderForm/CustomerForm";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { useTranslatedError } from "@/shared/hooks/useTranslatedError";
import { zodResolver } from "@hookform/resolvers/zod/dist/zod.js";
import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import { FormProvider, useForm } from "react-hook-form";
import zod from "zod";

const updateCustomerSchema = zod.object({
  customerFirstName: zod.string().min(1, "שדה חובה"),
  customerLastName: zod.string().min(1, "שדה חובה"),
  customerPhone: zod.string().min(1, "שדה חובה"),
  customerEmail: zod.email("כתובת אימייל לא חוקית").min(1, "שדה חובה"),
});
type UpdateCustomer = zod.infer<typeof updateCustomerSchema>;
export function ProfilePageContent() {
  const user = useAuthStore((state) => state.user);

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

  return (
    <main className="lg:w-2/3 lg:mx-auto">
      <h5 className="type-h5 text-navy my-6">הפרטים שלי</h5>
      <form onSubmit={handleSubmit((data) => mutate(data))}>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <FormProvider {...formMethods}>
            <CustomerForm />
          </FormProvider>
        </div>
        <Button variant="brand" loading={isPending} type="submit" className="my-5 font-bold px-10 py-5 rounded-md">
          שמירה
        </Button>
        {isSuccess && <SuccessBadge>הפרטים נשמרו בהצלחה</SuccessBadge>}
        <ErrorDisplay>{translatedError}</ErrorDisplay>
      </form>
    </main>
  );
}
