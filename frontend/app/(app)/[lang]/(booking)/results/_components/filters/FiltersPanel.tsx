import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import { CheckboxFilter } from "./CheckboxFilter";
import { useFilterOptions } from "../../_hooks/useFiltersOptions";
import type { SelectedFilters } from "../../_hooks/useCheckboxFilters";
import type { FilterConfig } from "../../../_components/_constants/filtersList";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

interface FiltersPanelProps {
  cars: booking.AvailableVehicle[];
  selectedFilters: SelectedFilters;
  onToggle: (filterId: FilterConfig["id"], value: string) => void;
  onClear: () => void;
  hasActiveFilters: boolean;
}

export function FiltersPanel({
  cars,
  selectedFilters,
  onToggle,
  onClear,
  hasActiveFilters,
}: FiltersPanelProps) {
  const t = useTranslations();
  const filtersOptions = useFilterOptions(cars);

  const visibleFilters = filtersOptions.filter(
    (filter) => filter.options.length > 1,
  );

  if (visibleFilters.length === 0) {
    return null;
  }

  return (
    <aside className="py-6">
      {visibleFilters.map((filter, index) => (
        <CheckboxFilter
          key={filter.id}
          title={t(filter.titleKey)}
          icon={filter.icon}
          selectedValues={selectedFilters.get(filter.id) ?? new Set()}
          onChange={(value) => onToggle(filter.id, value)}
          showDivider={index > 0}
          options={filter.options.map((value) => ({
            value,
            label: filter.getOptionLabel(value, t),
          }))}
        />
      ))}

      {hasActiveFilters ? (
        <Button
          variant="outline"
          type="button"
          onClick={onClear}
          className="type-paragraph font-normal flex items-center gap-1 mt-6"
        >
          <X />
          {t("booking.results.clearFilters")}
        </Button>
      ) : null}
    </aside>
  );
}
