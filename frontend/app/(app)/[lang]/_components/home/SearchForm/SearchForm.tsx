"use client";
import z from "zod";
import { useState } from "react";
import { useTranslations } from "next-intl";
import { LocationCombobox } from "./LocationCombobox";
import { Button } from "@/components/ui/button";
import { CalendarInput } from "./CalendarInput";
import { TimeSelect } from "./TimeSelect";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { DifferentLocCheckbox } from "./DifferentLocCheckbox";
import { AgePopover } from "./AgePopover";
import { CouponPopover } from "./CouponPopover";

const searchFormSchema = z.object({
  pickupLocation: z.int().min(1, "Pickup location is required"),
  dropoffLocation: z.int().min(1).optional(),
  pickupDate: z.date().min(new Date(), "Pickup date must be in the future"),
  pickupTime: z.string().min(1, "Pickup time is required"),
  dropoffDate: z
    .date()
    .min(new Date(), "Dropoff date must be in the future")
    .optional(),
  dropoffTime: z.string().min(1, "Dropoff time is required").optional(),
  driverAge: z.number().min(18, "Driver must be at least 18 years old"),
  couponCode: z.string().optional(),
});

export function SearchForm() {
  const t = useTranslations("SearchForm");
  const [isReturnDifferentLoc, setIsReturnDifferentLoc] = useState(false);

  const { control, handleSubmit } = useForm<z.infer<typeof searchFormSchema>>({
    resolver: zodResolver(searchFormSchema),
    defaultValues: {
      driverAge: 30,
      pickupTime: "",
      dropoffTime: "",
    },
  });

  function onSubmit(data: z.infer<typeof searchFormSchema>) {
    console.log("Form submitted with data:", data);
  }

  return (
    <form
      className="flex flex-col w-10/12 mx-auto mt-4"
      onSubmit={handleSubmit(onSubmit)}
    >
      <div className="bg-navy w-fit py-2 rounded-t-xl flex items-center text-white type-h6 px-6 gap-5">
        <DifferentLocCheckbox
          label={t("returnDifferentLoc")}
          isReturnDifferentLoc={isReturnDifferentLoc}
          setIsReturnDifferentLoc={setIsReturnDifferentLoc}
        />
        <div className="h-4 w-px bg-white/40 shrink-0" />
        <AgePopover
          checkboxLabel={t("ageRange")}
          inputLabel={t("agePopoverLabel")}
          saveButtonText={t("save")}
        />
        <div className="h-4 w-px bg-white/40 shrink-0" />
        <CouponPopover
          checkboxLabel={t("hasCoupon")}
          inputLabel={t("couponPlaceholder")}
          saveButtonText={t("save")}
        />
      </div>
      <div className="bg-white/95 w-full rounded-l-xl rounded-br-xl flex items-center gap-2 px-5">
        <div className="flex gap-2 flex-1 my-5 *:flex-1">
          <Controller
            name="pickupLocation"
            control={control}
            render={({ field }) => (
              <LocationCombobox
                placeholder={t("pickupLocationPlaceholder")}
                onSelect={(id) => field.onChange(id)}
              />
            )}
          />
          {isReturnDifferentLoc && (
            <Controller
              name="dropoffLocation"
              control={control}
              render={({ field }) => (
                <LocationCombobox
                  placeholder={t("dropoffLocationPlaceholder")}
                  onSelect={(id) => field.onChange(id)}
                />
              )}
            />
          )}
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <Controller
            name="pickupDate"
            control={control}
            render={({ field }) => (
              <CalendarInput
                placeholder={t("pickupDatePlaceholder")}
                value={field.value}
                onSelect={field.onChange}
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <Controller
            name="pickupTime"
            control={control}
            render={({ field }) => (
              <TimeSelect
                placeholder={t("timePlaceholder")}
                value={field.value}
                onChange={field.onChange}
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/10" : "w-1/6"}>
          <Controller
            name="dropoffDate"
            control={control}
            render={({ field }) => (
              <CalendarInput
                placeholder={t("dropoffDatePlaceholder")}
                value={field.value}
                onSelect={field.onChange}
              />
            )}
          />
        </div>
        <div className={isReturnDifferentLoc ? "w-1/12" : "w-1/10"}>
          <Controller
            name="dropoffTime"
            control={control}
            render={({ field }) => (
              <TimeSelect
                placeholder={t("timePlaceholder")}
                value={field.value}
                onChange={field.onChange}
              />
            )}
          />
        </div>
        <div className="w-1/9">
          <Button
            type="submit"
            variant="brand"
            className="w-full py-6.5 type-paragraph font-bold"
          >
            {t("searchButton")}
          </Button>
        </div>
      </div>
    </form>
  );
}
