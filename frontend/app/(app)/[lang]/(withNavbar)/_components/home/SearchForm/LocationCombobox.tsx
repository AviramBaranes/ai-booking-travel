import { useState } from "react";
import { useLocations } from "./useLocations";
import {
  Combobox,
  ComboboxContent,
  ComboboxInput,
  ComboboxList,
  ComboboxEmpty,
  ComboboxItem,
} from "@/components/ui/combobox";
import { Building2, MapPin, Plane } from "lucide-react";
import { booking } from "@/shared/client";
import { FieldError } from "react-hook-form";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";

interface LocationComboboxProps {
  placeholder: string;
  error: FieldError | undefined;
  onSelect: (locationId: number) => void;
  ref?: React.Ref<HTMLInputElement | null>;
  initializedLocations?: { id: number; name: string }[];
  value?: string;
}
export function LocationCombobox({
  placeholder,
  onSelect,
  error,
  ref,
  value,
  initializedLocations,
}: LocationComboboxProps) {
  const [search, setSearch] = useState("");
  const [selectedName, setSelectedName] = useState(value ?? "");
  const { locations } = useLocations(search);

  return (
    <Combobox
      items={locations?.length ? locations : initializedLocations || []}
      value={selectedName}
      onValueChange={(val) => {
        setSelectedName(val ?? "");
        const loc = locations.find((l) => l.name === val);
        if (loc) {
          onSelect(loc.id);
        }
      }}
    >
      <div className="flex flex-col">
        <ComboboxInput
          showClear={!!selectedName}
          placeholder={placeholder}
          aria-invalid={error ? "true" : "false"}
          className="search-form-input"
          showTrigger={false}
          onChange={(e) => setSearch(e.target.value)}
          readOnly={!!selectedName}
          ref={ref}
        >
          <MapPin className="absolute top-1/2 -translate-y-1/2 inset-s-3 size-4.5 text-brand pointer-events-none" />
        </ComboboxInput>
        <ErrorDisplay>{error?.message}</ErrorDisplay>
      </div>
      <ComboboxContent className="rounded-xl p-4">
        <ComboboxEmpty>לא נמצאו מיקומים</ComboboxEmpty>
        <ComboboxList className="divide-y divide-border" dir="ltr">
          {(loc: booking.LocationResult) => (
            <ComboboxItem
              key={loc.id}
              value={loc.name}
              className="flex items-center gap-3 px-3 py-3 text-base text-[#1b1b1b] rounded-none pr-3 pl-3 data-highlighted:text-brand data-highlighted:bg-[#f0f3f9]"
            >
              {loc.iata ? (
                <Plane className="size-5 shrink-0 text-brand!" />
              ) : (
                <Building2 className="size-5 shrink-0 text-brand!" />
              )}
              <span className="flex-1">{loc.name}</span>
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
