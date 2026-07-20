"use client";
import { useTranslations } from "next-intl";
import { LocationCombobox } from "./LocationCombobox";
import { Button } from "@/components/ui/button";
import { CalendarInput } from "./CalendarInput";
import { TimeSelect } from "./TimeSelect";
import { Controller, useWatch, useFormContext } from "react-hook-form";
import { DifferentLocCheckbox } from "./DifferentLocCheckbox";
import { AgePopover } from "./AgePopover";
import { CouponPopover } from "./CouponPopover";
import { SearchFormValues } from "./searchFormSchema";
import { useRef } from "react";
import clsx from "clsx";
import { CalendarInputRange } from "./CalendarInputRange";
import useAuthStore from "@/shared/auth/authStore";

export type SearchFieldHandle = {
  focus: () => void;
};

interface Location {
  id: number;
  name: string;
}

export interface SearchFormFields {
  loading: boolean;
  pickUpLocation: Location;
  dropOffLocation?: Location;
  pickUpDate: Date;
  dropOffDate: Date;
  pickUpTime: string;
  dropOffTime: string;
  driverAge: number;
  couponCode?: string;
}

interface SearchFormDesktopProps extends Partial<SearchFormFields> {
  className?: string;
}

export function SearchFormDesktop({
  className,
  loading,
  ...fields
}: SearchFormDesktopProps) {
  const t = useTranslations("SearchForm");

  const dropoffLocationRef = useRef<HTMLInputElement | null>(null);
  const pickupDateRef = useRef<SearchFieldHandle>(null);
  const dropoffDateRef = useRef<SearchFieldHandle>(null);
  const pickupTimeRef = useRef<SearchFieldHandle>(null);
  const dropoffTimeRef = useRef<SearchFieldHandle>(null);

  const { control, setValue } = useFormContext<SearchFormValues>();

  const isDropoffDifferentLoc =
    useWatch({
      control,
      name: "isDropoffDifferentLoc",
    }) ?? false;

  const pickupDate = useWatch({
    control,
    name: "pickupDate",
  });

  const user = useAuthStore((s) => s.user);
  const isAgent = user?.role === "agent";

  return (
    <>
      <div className="hidden lg:flex bg-navy w-fit py-2 rounded-t-xl items-center text-white type-h6 px-6 gap-5">
        <Controller
          name="isDropoffDifferentLoc"
          control={control}
          render={({ field }) => (
            <DifferentLocCheckbox
              label={t("dropoffDifferentLoc")}
              isDropoffDifferentLoc={field.value ?? false}
              setIsDropoffDifferentLoc={field.onChange}
            />
          )}
        />
        <div className="h-4 w-px bg-white/40 shrink-0" />
        <Controller
          name="driverAge"
          control={control}
          render={({ field }) => (
            <AgePopover
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
      <div className="flex flex-col lg:flex-row items-start justify-center gap-2 bg-white/95 w-full py-6 rounded-xl lg:max-h-35 min-h-25 lg:rounded-tr-none px-5">
        <div className="flex flex-col lg:flex-row gap-2 flex-1 *:flex-1 w-full">
          <Controller
            name="pickupLocation"
            control={control}
            render={({ field, fieldState }) => (
              <LocationCombobox
                placeholder={t("pickupLocationPlaceholder")}
                emptyMessage={t("noLocationsFound")}
                onSelect={(id) => {
                  field.onChange(id);
                  if (isDropoffDifferentLoc) {
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
          {isDropoffDifferentLoc && (
            <Controller
              name="dropoffLocation"
              control={control}
              render={({ field, fieldState }) => (
                <LocationCombobox
                  placeholder={t("dropoffLocationPlaceholder")}
                  emptyMessage={t("noLocationsFound")}
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

        <div className="flex w-full gap-2 lg:contents">
          <div
            className={clsx("w-1/2", {
              "lg:w-1/10": isDropoffDifferentLoc,
              "lg:w-1/6": !isDropoffDifferentLoc,
            })}
          >
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
                  disabledFn={(date) => {
                    const yesterday = new Date();
                    yesterday.setDate(yesterday.getDate() - 1);
                    return (
                      date < yesterday ||
                      date > new Date(Date.now() + 365 * 24 * 60 * 60 * 1000)
                    );
                  }}
                />
              )}
            />
          </div>
          <div
            className={clsx("w-1/2", {
              "lg:w-1/12": isDropoffDifferentLoc,
              "lg:w-1/10": !isDropoffDifferentLoc,
            })}
          >
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

                    requestAnimationFrame(() => {
                      requestAnimationFrame(() => {
                        dropoffDateRef.current?.focus();
                      });
                    });
                  }}
                  error={fieldState.error}
                />
              )}
            />
          </div>
        </div>
        <div className="flex w-full gap-2 lg:contents">
          <div
            className={clsx("w-1/2", {
              "lg:w-1/10": isDropoffDifferentLoc,
              "lg:w-1/6": !isDropoffDifferentLoc,
            })}
          >
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
          <div
            className={clsx("w-1/2", {
              "lg:w-1/12": isDropoffDifferentLoc,
              "lg:w-1/10": !isDropoffDifferentLoc,
            })}
          >
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
        </div>
        <div className="w-full lg:w-1/9">
          <Button
            type="submit"
            variant="brand"
            className="w-full py-8 lg:py-6 type-paragraph font-bold cursor-pointer"
            loading={loading}
          >
            {t("searchButton")}
          </Button>
        </div>
      </div>
    </>
  );
}
