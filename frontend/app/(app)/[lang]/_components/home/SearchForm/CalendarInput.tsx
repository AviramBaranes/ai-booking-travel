import { useState } from "react";
import { he } from "react-day-picker/locale";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { Calendar as CalendarIcon } from "lucide-react";
import { Field } from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { useParams } from "next/navigation";

interface CalendarInputProps {
  placeholder: string;
}

export function CalendarInput({ placeholder }: CalendarInputProps) {
  const { lang } = useParams();

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Field>
          <InputGroup className="search-form-input px-0">
            <InputGroupInput
              id="input-group-url"
              placeholder={placeholder}
              className="text-start px-2"
            />
            <InputGroupAddon align="inline-start" className="pl-1 pr-0">
              <CalendarIcon className="size-5 mr-2 text-brand" />
            </InputGroupAddon>
          </InputGroup>
        </Field>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          locale={lang === "he" ? he : undefined}
          numberOfMonths={2}
        />
      </PopoverContent>
    </Popover>
  );
}
