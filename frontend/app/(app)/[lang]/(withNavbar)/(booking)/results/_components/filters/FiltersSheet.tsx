import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { X, SlidersVertical } from "lucide-react";
import { useTranslations } from "next-intl";
import { CarGroupsFilter } from "./CarGroupsFilter";
import { availability } from "@/shared/client";
import { FiltersPanel } from "./FiltersPanel";
import { Button } from "@/components/ui/button";

interface FiltersSheetProps {
  cars: availability.AvailableVehicle[];
  hasActiveFilters: boolean;
}
export function FiltersSheet({ cars, hasActiveFilters }: FiltersSheetProps) {
  const t = useTranslations("booking.results.filters.mobileSheet");

  return (
    <Sheet>
      <SheetTrigger asChild>
        <div className="lg:hidden fixed bottom-0 w-full bg-white p-4 shadow-card flex justify-center z-50">
          <Button variant="ghost">
            <SlidersVertical className="mr-2 text-brand size-5 shrink-0" />
            <span className="type-paragraph font-semibold text-navy">
              {t("buttonTitle")}
            </span>
          </Button>
        </div>
      </SheetTrigger>
      <SheetContent
        side="top"
        showCloseButton={false}
        className="p-0 rounded-none border-0 w-full bottom-0 flex flex-col lg:hidden"
      >
        <div className="flex items-center justify-between mt-12 ">
          <SheetTitle className="mx-5 type-h5 text-navy">{t("title")}</SheetTitle>
          <div className="flex items-center justify-between px-4 mx-10 shrink-0">
            <SheetClose asChild>
              <button aria-label="Close menu">
                <X className="size-5 text-navy" />
              </button>
            </SheetClose>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto">
          <div className="mx-5">
            <CarGroupsFilter title={""} />
          </div>
          <div className="mx-11">
            <FiltersPanel cars={cars} hasActiveFilters={hasActiveFilters} />
          </div>
        </div>

        {/* Apply button fixed to the bottom */}
        <div className="px-5 pt-4 pb-[max(1rem,env(safe-area-inset-bottom))] border-t border-border-light shrink-0">
          <SheetClose asChild>
            <Button
              type="button"
              variant="brand"
              className="w-full py-6 type-paragraph font-bold cursor-pointer"
            >
              {t("applyButton")}
            </Button>
          </SheetClose>
        </div>
      </SheetContent>
    </Sheet>
  );
}
