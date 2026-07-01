"use client";

import { useRef, useState } from "react";

import {
  Sheet,
  SheetContent,
  SheetClose,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { LocationCombobox } from "./LocationCombobox";
import { FieldError } from "react-hook-form";
import { X } from "lucide-react";

interface LocationComboboxSheetProps {
  placeholder: string;
  error?: FieldError;
  onSelect: (locationId: number) => void;
  initializedLocations?: { id: number; name: string }[];
  value?: string;
}

export function LocationComboboxSheet({
  placeholder,
  error,
  onSelect,
  initializedLocations,
  value,
}: LocationComboboxSheetProps) {
  const ref = useRef<HTMLInputElement | null>(null);
  const [open, setOpen] = useState(false);
  const [selectedName, setSelectedName] = useState(value ?? "");
  const [sheetContentEl, setSheetContentEl] = useState<HTMLElement | null>(
    null,
  );

  const close = () => setOpen(false);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <div
          onClick={(e) => {
            e.stopPropagation();
            setOpen(true);
            requestAnimationFrame(() => {
              ref.current?.focus();
            });
          }}
        >
          <LocationCombobox
            placeholder={placeholder}
            error={error}
            onSelect={() => {}}
            value={selectedName}
            open={false}
          />
        </div>
      </SheetTrigger>

      <SheetContent
        side="top"
        ref={setSheetContentEl}
        showCloseButton={false}
        className="p-0 rounded-none border-0 w-full bottom-0"
      >
        <SheetTitle className="sr-only">Location Search</SheetTitle>
        {/* Top bar */}
        <div className="flex h-15 items-center justify-between px-4 border-b border-border-light">
          <SheetClose asChild>
            <button aria-label="Close menu">
              <X className="size-5 text-navy" />
            </button>
          </SheetClose>
        </div>
        <div className="w-11/12 mx-auto">
          <LocationCombobox
            ref={ref}
            placeholder={placeholder}
            error={error}
            onSelect={(locationId, name) => {
              onSelect(locationId);
              setSelectedName(name);
              close();
            }}
            initializedLocations={initializedLocations}
            value={selectedName}
            container={sheetContentEl}
          />
        </div>
      </SheetContent>
    </Sheet>
  );
}
