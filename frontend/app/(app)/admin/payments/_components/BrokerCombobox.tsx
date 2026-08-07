"use client";

import {
  Combobox,
  ComboboxInput,
  ComboboxContent,
  ComboboxList,
  ComboboxItem,
  ComboboxEmpty,
  ComboboxCollection,
} from "@/components/ui/combobox";

export type Broker = {
  code: string;
  name: string;
};

// Only Flex is settled through this screen for now; Hertz joins the list once its
// payment flow exists.
export const BROKERS: Broker[] = [{ code: "flex", name: "פלקס" }];

interface BrokerComboboxProps {
  value: Broker | null;
  onChange: (value: Broker | null) => void;
}

export function BrokerCombobox({ value, onChange }: BrokerComboboxProps) {
  return (
    <Combobox
      items={BROKERS}
      value={value}
      onValueChange={onChange}
      itemToStringLabel={(item: Broker) => item.name}
      itemToStringValue={(item: Broker) => item.code}
      isItemEqualToValue={(a: Broker, b: Broker) => a.code === b.code}
    >
      <ComboboxInput placeholder="בחר ספק..." showClear className="w-full" />
      <ComboboxContent>
        <ComboboxEmpty>לא נמצאו תוצאות</ComboboxEmpty>
        <ComboboxList>
          <ComboboxCollection>
            {(item: Broker) => (
              <ComboboxItem key={item.code} value={item}>
                {item.name}
              </ComboboxItem>
            )}
          </ComboboxCollection>
        </ComboboxList>
      </ComboboxContent>
    </Combobox>
  );
}
