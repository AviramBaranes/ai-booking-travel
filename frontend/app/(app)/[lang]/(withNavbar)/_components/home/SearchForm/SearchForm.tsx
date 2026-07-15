"use client";
import { useTranslations } from "next-intl";
import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { SearchFormValues, searchSchema } from "./searchFormSchema";
import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import clsx from "clsx";
import { formatDate } from "@/shared/utils/formatDate";
import { SearchFormDesktop } from "./SearchFormDesktop";
import { SearchFormMobile } from "./SearchFormMobile";
import useAuthStore from "@/shared/auth/authStore";
import { OPEN_DIALOG_QUERY_KEY, OPEN_DIALOG_QUERY_VALUE } from "../../header/login/useDialogOpenFromQuery";

export type SearchFieldHandle = {
  focus: () => void;
};

interface Location {
  id: number;
  name: string;
}

export interface SearchFormFields {
  pickUpLocation: Location;
  dropOffLocation?: Location;
  pickUpDate: Date;
  dropOffDate: Date;
  pickUpTime: string;
  dropOffTime: string;
  driverAge: number;
  couponCode?: string;
}

interface SearchFormProps extends Partial<SearchFormFields> {
  className?: string;
}

export function SearchForm({ className, ...fields }: SearchFormProps) {
  const { lang } = useParams();
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const clearSession = useBookingSessionStore((s) => s.clearSession);
  const clearCarGroupFilters = useBookingSessionStore(
    (s) => s.clearCarGroupFilters,
  );
  const clearAllCheckboxFilters = useBookingSessionStore(
    (s) => s.clearAllCheckboxFilters,
  );
  const t = useTranslations("SearchForm");
  const searchFormSchema = searchSchema(t);
  const [loading, setLoading] = useState(false);

  const formMethods = useForm<SearchFormValues>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: {
      isDropoffDifferentLoc:
        !!fields.dropOffLocation &&
        fields.dropOffLocation.id !== fields.pickUpLocation?.id,
      driverAge: fields.driverAge ?? 30,
      pickupTime: fields.pickUpTime ?? "",
      dropoffTime: fields.dropOffTime ?? "",
      couponCode: fields.couponCode ?? "",
      pickupLocation: fields.pickUpLocation?.id,
      dropoffLocation: fields.dropOffLocation?.id,
      pickupDate: fields.pickUpDate ?? undefined,
      dropoffDate: fields.dropOffDate ?? undefined,
    },
  });

  function onSubmit(data: SearchFormValues) {
    clearSession();
    clearCarGroupFilters();
    clearAllCheckboxFilters();
    setLoading(true);
    const urlParams = new URLSearchParams();

    urlParams.set("pl", data.pickupLocation.toString());
    urlParams.set(
      "dl",
      data.isDropoffDifferentLoc
        ? data.dropoffLocation!.toString()
        : data.pickupLocation.toString(),
    );
    urlParams.set("pd", formatDate(data.pickupDate!));
    urlParams.set("pt", data.pickupTime);
    urlParams.set("dd", formatDate(data.dropoffDate!));
    urlParams.set("dt", data.dropoffTime);
    urlParams.set("da", data.driverAge.toString());

    if (data.couponCode) {
      urlParams.set("cc", data.couponCode);
    }

    setTimeout(() => {
      setLoading(false);
    }, 2000);

    location.href = `/${lang}/results?${urlParams.toString()}`;
  }

  const isAgent = user?.role === "agent";

  return (
    <form
      className={clsx("flex flex-col w-10/12 mx-auto mt-4", className)}
      onSubmit={formMethods.handleSubmit(onSubmit)}
      onClick={(e) => {
        if (!isAgent) {
          e.stopPropagation();
          router.push(
            `/${lang}?${OPEN_DIALOG_QUERY_KEY}=${OPEN_DIALOG_QUERY_VALUE}`,
          );
        }
      }}
    >
      <FormProvider {...formMethods}>
        <div className="hidden lg:block">
          <SearchFormDesktop
            loading={loading}
            className={className}
            {...fields}
          />
        </div>
        <div className="block lg:hidden">
          <SearchFormMobile
            loading={loading}
            className={className}
            {...fields}
          />
        </div>
      </FormProvider>
    </form>
  );
}
