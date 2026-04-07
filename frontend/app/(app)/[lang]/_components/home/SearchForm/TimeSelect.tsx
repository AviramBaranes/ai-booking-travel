"use client";

import * as React from "react";
import { Check, Clock } from "lucide-react";
import { Field } from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const times = Array.from(
  { length: 48 },
  (_, i) =>
    `${String((6 + Math.floor(i / 2)) % 24).padStart(2, "0")}:${i % 2 ? "30" : "00"}`,
);

interface TimeSelectProps {
  placeholder: string;
  value?: string;
  onChange?: (value: string) => void;
  name?: string;
}

export function TimeSelect({
  placeholder,
  value,
  onChange,
  name,
}: TimeSelectProps) {
  const [internalValue, setInternalValue] = React.useState("");

  const currentValue = value ?? internalValue;

  const handleSelect = (nextValue: string) => {
    if (value === undefined) {
      setInternalValue(nextValue);
    }

    onChange?.(nextValue);
  };

  return (
    <div className="w-full">
      {name ? <input type="hidden" name={name} value={currentValue} /> : null}

      <DropdownMenu modal={false}>
        <DropdownMenuTrigger asChild>
          <div className="w-full cursor-pointer">
            <Field>
              <InputGroup className="search-form-input px-0">
                <InputGroupInput
                  value={currentValue}
                  placeholder={placeholder}
                  readOnly
                  className="text-start px-2 cursor-pointer"
                />
                <InputGroupAddon align="inline-start" className="pl-1 pr-0">
                  <Clock className="size-5 mr-2 text-brand" />
                </InputGroupAddon>
              </InputGroup>
            </Field>
          </div>
        </DropdownMenuTrigger>

        <DropdownMenuContent
          align="start"
          sideOffset={8}
          className="w-(--radix-dropdown-menu-trigger-width) rounded-xl p-2"
        >
          <div className="max-h-72 overflow-y-auto">
            {times.map((time) => {
              const selected = currentValue === time;

              return (
                <DropdownMenuItem
                  key={time}
                  onClick={() => handleSelect(time)}
                  className="flex items-center justify-between rounded-md px-3 py-3 text-base"
                >
                  <span>{time}</span>
                  <Check
                    className={
                      selected ? "size-4 text-brand" : "size-4 opacity-0"
                    }
                  />
                </DropdownMenuItem>
              );
            })}
          </div>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}
