import { useEffect, useState } from "react";
import { useLocations } from "./useLocations";
import {
  Combobox,
  ComboboxContent,
  ComboboxInput,
  ComboboxList,
  ComboboxEmpty,
  ComboboxItem,
  useComboboxAnchor,
} from "@/components/ui/combobox";
import { Building2, MapPin, Plane } from "lucide-react";
import { location } from "@/shared/client";
import { FieldError } from "react-hook-form";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { useDirection } from "@/shared/hooks/useDirection";
import { clsx } from "clsx";

interface LocationComboboxProps {
  placeholder: string;
  emptyMessage: string;
  error: FieldError | undefined;
  onSelect: (locationId: number, name: string) => void;
  ref?: React.Ref<HTMLInputElement | null>;
  initializedLocations?: { id: number; name: string }[];
  value?: string;
  open?: boolean;
  container?: HTMLElement | null;
}
export function LocationCombobox({
  placeholder,
  emptyMessage,
  onSelect,
  error,
  ref,
  value,
  initializedLocations,
  open,
  container,
}: LocationComboboxProps) {
  const dir = useDirection();
  const [search, setSearch] = useState("");
  const anchorRef = useComboboxAnchor();
  const [selectedName, setSelectedName] = useState(value ?? "");
  const { locations } = useLocations(search);

  useEffect(() => {
    setSelectedName(value ?? "");
  }, [value]);

  const items = locations?.length ? locations : initializedLocations || [];

  return (
    <Combobox
      items={items}
      filteredItems={items}
      value={selectedName}
      open={open}
      onValueChange={(val) => {
        setSelectedName(val ?? "");
        const loc = locations.find((l) => l.name === val);
        if (loc) {
          onSelect(loc.id, loc.name);
        }
      }}
    >
      <div ref={anchorRef} className="flex w-full flex-col">
        <ComboboxInput
          showClear={!!selectedName}
          placeholder={placeholder}
          aria-invalid={error ? "true" : "false"}
          inputClassName="text-sm"
          className="search-form-input md:text-base px-7"
          clearClassName={clsx("p-0 absolute", {
            "left-3": dir === "rtl",
            "right-3": dir === "ltr",
          })}
          showTrigger={false}
          onChange={(e) => setSearch(e.target.value)}
          readOnly={!!selectedName}
          ref={ref}
        >
          <MapPin className="pointer-events-none absolute inset-s-3 top-1/2 size-4.5 -translate-y-1/2 text-brand" />
        </ComboboxInput>
        <ErrorDisplay>{error?.message}</ErrorDisplay>
      </div>
      <ComboboxContent
        anchor={anchorRef}
        container={container}
        align="start"
        className="w-(--anchor-width)! min-w-(--anchor-width)! max-w-(--anchor-width)! rounded-xl p-1"
      >
        <ComboboxEmpty>{emptyMessage}</ComboboxEmpty>
        <ComboboxList className="divide-y divide-border" dir="ltr">
          {(loc: location.LocationResult) => (
            <ComboboxItem
              key={loc.id}
              value={loc.name}
              className="flex items-center gap-3 px-3 py-3 text-base text-[#1b1b1b] rounded-none pr-3 pl-3 data-highlighted:text-brand data-highlighted:bg-[#f0f3f9]"
            >
              {loc.isAirport ? (
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
