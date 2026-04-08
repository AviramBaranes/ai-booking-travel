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

export type SearchFieldHandle = {
  focus: () => void;
};

export function SearchForm() {
  const router = useRouter();
  const { lang } = useParams();
  const t = useTranslations("SearchForm");
  const searchFormSchema = searchSchema(t);

  const dropoffLocationRef = useRef<HTMLInputElement | null>(null);
  const pickupDateRef = useRef<SearchFieldHandle>(null);
  const dropoffDateRef = useRef<SearchFieldHandle>(null);
  const pickupTimeRef = useRef<SearchFieldHandle>(null);
  const dropoffTimeRef = useRef<SearchFieldHandle>(null);

  const { control, handleSubmit } = useForm<SearchFormValues>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: {
      isReturnDifferentLoc: false,
      driverAge: 30,
      pickupTime: "",
      dropoffTime: "",
    },
  });

  const isReturnDifferentLoc =
    useWatch({
      control,
      name: "isReturnDifferentLoc",
    }) ?? false;

  function formatDate(date: Date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  }

  function onSubmit(data: SearchFormValues) {
    const urlParams = new URLSearchParams();

    urlParams.set("pickupLocation", data.pickupLocation.toString());
    urlParams.set(
      "dropoffLocation",
      data.isReturnDifferentLoc
        ? data.dropoffLocation!.toString()
        : data.pickupLocation.toString(),
    );
    urlParams.set("pickupDate", formatDate(data.pickupDate!));
    urlParams.set("pickupTime", data.pickupTime);
    urlParams.set("dropoffDate", formatDate(data.dropoffDate!));
    urlParams.set("dropoffTime", data.dropoffTime);
    urlParams.set("driverAge", data.driverAge.toString());

    if (data.couponCode) {
      urlParams.set("couponCode", data.couponCode);
    }

    router.push(`/${lang}/results?${urlParams.toString()}`);
  }

  return (
    <form
      className="flex flex-col w-10/12 mx-auto mt-4"
      onSubmit={handleSubmit(onSubmit)}
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
                  pickupTimeRef.current?.focus();
                }}
                error={fieldState.error}
                ref={pickupDateRef}
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
              <CalendarInput
                ref={dropoffDateRef}
                placeholder={t("dropoffDatePlaceholder")}
                value={field.value}
                onSelect={(e) => {
                  field.onChange(e);
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
            className="w-full py-6 type-paragraph font-bold"
          >
            {t("searchButton")}
          </Button>
        </div>
      </div>
    </form>
  );
}
