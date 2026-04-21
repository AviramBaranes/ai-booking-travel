"use client";
import { useTranslations } from "next-intl";
import { LocationCombobox } from "./LocationCombobox";
import { Button } from "@/components/ui/button";
import { CalendarInput } from "./CalendarInput";
import { TimeSelect } from "./TimeSelect";
import { useForm, Controller, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { DifferentLocCheckbox } from "./DifferentLocCheckbox";
import { AgePopover } from "./AgePopover";
import { CouponPopover } from "./CouponPopover";
import { SearchFormValues, searchSchema } from "./searchFormSchema";
import { useRef } from "react";
import { useParams, useRouter } from "next/navigation";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";
import clsx from "clsx";
import { CalendarInputRange } from "./CalendarInputRange";
import { useSession } from "next-auth/react";
import {
  OPEN_DIALOG_QUERY_KEY,
  OPEN_DIALOG_QUERY_VALUE,
} from "../../header/login/useDialogOpenFromQuery";

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
  const session = useSession();
  const isAuthenticated = session.status === "authenticated";
  const isAgent = session.data?.user?.role === "agent";
  const router = useRouter();
  const { lang } = useParams();
  const clearSession = useBookingSessionStore((s) => s.clearSession);
  const t = useTranslations("SearchForm");
  const searchFormSchema = searchSchema(t);

  const dropoffLocationRef = useRef<HTMLInputElement | null>(null);
  const pickupDateRef = useRef<SearchFieldHandle>(null);
  const dropoffDateRef = useRef<SearchFieldHandle>(null);
  const pickupTimeRef = useRef<SearchFieldHandle>(null);
  const dropoffTimeRef = useRef<SearchFieldHandle>(null);

  const { control, handleSubmit, setValue } = useForm<SearchFormValues>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: {
      isReturnDifferentLoc:
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

  const isReturnDifferentLoc =
    useWatch({
      control,
      name: "isReturnDifferentLoc",
    }) ?? false;

  const pickupDate = useWatch({
    control,
    name: "pickupDate",
  });

  function formatDate(date: Date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }

  function onSubmit(data: SearchFormValues) {
    clearSession();
    const urlParams = new URLSearchParams();

    urlParams.set("pl", data.pickupLocation.toString());
    urlParams.set(
      "rl",
      data.isReturnDifferentLoc
        ? data.dropoffLocation!.toString()
        : data.pickupLocation.toString(),
    );
    urlParams.set("pd", formatDate(data.pickupDate!));
    urlParams.set("pt", data.pickupTime);
    urlParams.set("rd", formatDate(data.dropoffDate!));
    urlParams.set("rt", data.dropoffTime);
    urlParams.set("da", data.driverAge.toString());

    if (data.couponCode) {
      urlParams.set("cc", data.couponCode);
    }

    router.push(`/${lang}/results?${urlParams.toString()}`);
  }

  return (
    <form
      className={clsx("flex flex-col w-10/12 mx-auto mt-4", className)}
      onSubmit={handleSubmit(onSubmit)}
      onClick={(e) => {
        if (!isAuthenticated) {
          e.stopPropagation();
          router.push(
            `/${lang}?${OPEN_DIALOG_QUERY_KEY}=${OPEN_DIALOG_QUERY_VALUE}`,
          );
        }
      }}
    >
      <div className="bg-navy w-fit py-2 rounded-t-xl flex items-center text-white type-h6 px-6 gap-5">
        <Controller
          name="isReturnDifferentLoc"
          control={control}
          render={({ field }) => (
            <DifferentLocCheckbox
              label={t("returnDifferentLoc")}
              isReturnDifferentLoc={field.value ?? false}
              setIsReturnDifferentLoc={field.onChange}
            />
          )}
        />
        <div className="h-4 w-px bg-white/40 shrink-0" />
        <Controller
          name="driverAge"
          control={control}
          render={({ field }) => (
            <AgePopover
              checkboxLabel={t("ageRange")}
              inputLabel={t("agePopoverLabel")}
              saveButtonText={t("save")}
              driverAge={field.value}
              setDriverAge={field.onChange}
            />
          )}
        />
        {!isAgent && (
          <>
            <div className="h-4 w-px bg-white/40 shrink-0" />

            <Controller
              name="couponCode"
              control={control}
              render={({ field }) => (
                <CouponPopover
                  checkboxLabel={t("hasCoupon")}
                  inputLabel={t("couponPlaceholder")}
                  saveButtonText={t("save")}
                  couponCode={field.value ?? ""}
                  setCouponCode={field.onChange}
                />
              )}
            />
          </>
        )}
      </div>
      <div className="bg-white/95 w-full py-6 rounded-l-xl max-h-35 min-h-25 justify-center rounded-br-xl flex items-start gap-2 px-5">
        <div className="flex gap-2 flex-1 *:flex-1">
          <Controller
            name="pickupLocation"
            control={control}
            render={({ field, fieldState }) => (
              <LocationCombobox
                placeholder={t("pickupLocationPlaceholder")}
                onSelect={(id) => {
                  field.onChange(id);
                  if (isReturnDifferentLoc) {
                    dropoffLocationRef.current?.focus();
                  } else {
                    pickupDateRef.current?.focus();
                  }
                }}
                error={fieldState.error}
                value={fields.pickUpLocation?.name ?? ""}
                initializedLocations={
                  fields.pickUpLocation ? [fields.pickUpLocation] : undefined
                }
              />
            )}
          />
          {isReturnDifferentLoc && (
            <Controller
              name="dropoffLocation"
              control={control}
              render={({ field, fieldState }) => (
                <LocationCombobox
                  placeholder={t("dropoffLocationPlaceholder")}
                  onSelect={(id) => {
                    field.onChange(id);
                    pickupDateRef.current?.focus();
                  }}
                  error={fieldState.error}
                  ref={dropoffLocationRef}
                  value={fields.dropOffLocation?.name ?? ""}
                  initializedLocations={
                    fields.dropOffLocation
                      ? [fields.dropOffLocation]
                      : undefined
                  }
                />
              )}
            />
          )}
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <Controller
            name="pickupDate"
            control={control}
            render={({ field, fieldState }) => (
              <CalendarInput
                placeholder={t("pickupDatePlaceholder")}
                value={field.value}
                onSelect={(e) => {
                  field.onChange(e);
                  setValue("dropoffDate", null);
                  pickupTimeRef.current?.focus();
                }}
                error={fieldState.error}
                ref={pickupDateRef}
                disabledFn={(date) =>
                  date < new Date() ||
                  date > new Date(Date.now() + 365 * 24 * 60 * 60 * 1000)
                }
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <Controller
            name="pickupTime"
            control={control}
            render={({ field, fieldState }) => (
              <TimeSelect
                ref={pickupTimeRef}
                placeholder={t("timePlaceholder")}
                value={field.value}
                onChange={(e) => {
                  field.onChange(e);
                  dropoffDateRef.current?.focus();
                }}
                error={fieldState.error}
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <Controller
            name="dropoffDate"
            control={control}
            render={({ field, fieldState }) => (
              <CalendarInputRange
                ref={dropoffDateRef}
                placeholder={t("dropoffDatePlaceholder")}
                value={{ from: pickupDate, to: field.value }}
                onSelect={(e) => {
                  field.onChange(e?.to);
                  dropoffTimeRef.current?.focus();
                }}
                error={fieldState.error}
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <Controller
            name="dropoffTime"
            control={control}
            render={({ field, fieldState }) => (
              <TimeSelect
                ref={dropoffTimeRef}
                placeholder={t("timePlaceholder")}
                value={field.value ?? ""}
                onChange={field.onChange}
                error={fieldState.error}
              />
            )}
          />
        </div>
        <div className="w-1/9">
          <Button
            type="submit"
            variant="brand"
            className="w-full py-6 type-paragraph font-bold cursor-pointer"
          >
            {t("searchButton")}
          </Button>
        </div>
      </div>
    </form>
  );
}
