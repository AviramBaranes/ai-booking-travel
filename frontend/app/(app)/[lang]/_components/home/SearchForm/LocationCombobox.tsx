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

interface LocationComboboxProps {
  placeholder: string;
}
export function LocationCombobox({ placeholder }: LocationComboboxProps) {
  const [search, setSearch] = useState("");
  const { locations } = useLocations(search);

  return (
    <Combobox items={locations}>
      <ComboboxInput
        placeholder={placeholder}
        className="w-full rounded-md border-[1.5px] px-8 py-6 bg-input-bg border-input-border text-navy placeholder:text-secondary focus-within:border-brand-blue/50 [&_input]:bg-transparent"
        showTrigger={false}
        onChange={(e) => setSearch(e.target.value)}
      >
        <MapPin className="absolute top-1/2 -translate-y-1/2 inset-s-3 size-4.5 text-brand pointer-events-none" />
      </ComboboxInput>
      <ComboboxContent className="rounded-xl p-4">
        <ComboboxEmpty>לא נמצאו מיקומים</ComboboxEmpty>
        <ComboboxList className="divide-y divide-border" dir="ltr">
          {(loc: booking.LocationResult) => (
            <ComboboxItem
              key={loc.id}
              value={loc.name}
              className="flex items-center gap-3 px-3 py-3 text-base text-[#1b1b1b] rounded-none pr-3 pl-3 data-highlighted:bg-[#f0f3f9]"
            >
              {loc.iata ? (
                <Plane className="size-5 shrink-0 text-brand" />
              ) : (
                <Building2 className="size-5 shrink-0 text-brand" />
              )}
              <span className="flex-1">{loc.name}</span>
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
