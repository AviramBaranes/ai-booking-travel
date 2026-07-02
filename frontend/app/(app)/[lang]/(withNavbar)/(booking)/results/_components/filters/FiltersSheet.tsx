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
        className="p-0 rounded-none border-0 w-full bottom-0 flex flex-col lg:hidden overflow-y-scroll"
      >
        <div className="flex items-center justify-between mt-12 ">
          <SheetTitle className="mx-5">
            <h5 className="type-h5 text-navy">{t("title")}</h5>
          </SheetTitle>
          <div className="flex items-center justify-between px-4 mx-10 shrink-0">
            <SheetClose asChild>
              <button aria-label="Close menu">
                <X className="size-5 text-navy" />
              </button>
            </SheetClose>
          </div>
        </div>

        <div className="mx-5">
          <CarGroupsFilter title={""} />
        </div>
        <div className="mx-11">
          <FiltersPanel cars={cars} hasActiveFilters={hasActiveFilters} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
