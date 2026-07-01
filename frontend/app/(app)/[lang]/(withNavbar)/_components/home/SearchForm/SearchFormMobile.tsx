"use client";
import { useTranslations } from "next-intl";
import { LocationCombobox } from "./LocationCombobox";
import { Button } from "@/components/ui/button";
import { TimeSelect } from "./TimeSelect";
import { Controller, useWatch, useFormContext } from "react-hook-form";
import { DifferentLocCheckbox } from "./DifferentLocCheckbox";
import { AgePopover } from "./AgePopover";
import { CouponPopover } from "./CouponPopover";
import { SearchFormValues } from "./searchFormSchema";
import { useRef, useState } from "react";
import clsx from "clsx";
import { CalendarSheet, CalendarInputTrigger } from "./CalendarSheet";
import { LocationComboboxSheet } from "./LocationComboboxSheet";
import { AgeCheckbox } from "./AgeCheckbox";
import { CouponCheckbox } from "./CouponCheckbox";

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

interface SearchFormMobileProps extends Partial<SearchFormFields> {
  className?: string;
}

export function SearchFormMobile({
  className,
  loading,
  ...fields
}: SearchFormMobileProps) {
  const t = useTranslations("SearchForm");

  const dropoffLocationRef = useRef<HTMLInputElement | null>(null);
  const pickupTimeRef = useRef<SearchFieldHandle>(null);
  const dropoffTimeRef = useRef<SearchFieldHandle>(null);

  const [calendarOpen, setCalendarOpen] = useState(false);

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

  const dropoffDate = useWatch({
    control,
    name: "dropoffDate",
  });

  return (
    <div className="flex flex-col items-start justify-center gap-2 bg-white w-full rounded-xl p-2">
      <div className="flex flex-col gap-2 w-full">
        <Controller
          name="pickupLocation"
          control={control}
          render={({ field, fieldState }) => (
            <LocationComboboxSheet
              placeholder={t("pickupLocationPlaceholder")}
              onSelect={(id) => {
                field.onChange(id);
                if (isDropoffDifferentLoc) {
                  dropoffLocationRef.current?.focus();
                } else {
                  setCalendarOpen(true);
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
        {isDropoffDifferentLoc && (
          <Controller
            name="dropoffLocation"
            control={control}
            render={({ field, fieldState }) => (
              <LocationComboboxSheet
                placeholder={t("dropoffLocationPlaceholder")}
                onSelect={(id) => {
                  field.onChange(id);
                  setCalendarOpen(true);
                }}
                error={fieldState.error}
                value={fields.dropOffLocation?.name ?? ""}
                initializedLocations={
                  fields.dropOffLocation ? [fields.dropOffLocation] : undefined
                }
              />
            )}
          />
        )}
      </div>

      <div className="flex w-full gap-2 lg:contents">
        <div className="w-1/2">
          <Controller
            name="pickupDate"
            control={control}
            render={({ field, fieldState }) => (
              <CalendarInputTrigger
                placeholder={t("pickupDatePlaceholder")}
                value={field.value}
                error={fieldState.error}
                onClick={() => setCalendarOpen(true)}
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
                      dropoffTimeRef.current?.focus();
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
              <CalendarInputTrigger
                placeholder={t("dropoffDatePlaceholder")}
                value={field.value}
                error={fieldState.error}
                onClick={() => setCalendarOpen(true)}
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
      <Controller
        name="driverAge"
        control={control}
        render={({ field }) => (
          <AgeCheckbox
            saveButtonText={t("save")}
            driverAge={field.value}
            setDriverAge={field.onChange}
          />
        )}
      />
      <Controller
          name="couponCode"
          control={control}
          render={({ field }) => (
            <CouponCheckbox
              checkboxLabel={t("hasCoupon")}
              inputLabel={t("couponPlaceholder")}
              saveButtonText={t("save")}
              couponCode={field.value ?? ""}
              setCouponCode={field.onChange}
            />
          )}
        />
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

      <CalendarSheet
        open={calendarOpen}
        onOpenChange={setCalendarOpen}
        value={{ from: pickupDate, to: dropoffDate }}
        onConfirm={(range) => {
          setValue("pickupDate", range?.from ?? null, {
            shouldValidate: true,
          });
          setValue("dropoffDate", range?.to ?? null, {
            shouldValidate: true,
          });
          pickupTimeRef.current?.focus();
        }}
      />
    </div>
  );
}
