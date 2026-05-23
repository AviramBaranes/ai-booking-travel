"use client";

import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
} from "@/components/ui/combobox";

interface ReportSupplierComboboxProps {
  value: string;
  onChange: (value: string) => void;
}

const supplierOptions: string[] = [];

export function ReportSupplierCombobox({
  value,
  onChange,
}: ReportSupplierComboboxProps) {
  return (
    <Combobox
      items={supplierOptions}
      value={value || null}
      onValueChange={(nextValue: string | null) => onChange(nextValue ?? "")}
    >
      <ComboboxInput
        placeholder="כל הספקים"
        showClear
        className="w-full"
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
      <ComboboxContent>
        <ComboboxEmpty>אין ספקים להצגה כרגע</ComboboxEmpty>
        <ComboboxList>
          {(supplier) => (
            <ComboboxItem key={supplier} value={supplier}>
              {supplier}
            </ComboboxItem>
          )}
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
