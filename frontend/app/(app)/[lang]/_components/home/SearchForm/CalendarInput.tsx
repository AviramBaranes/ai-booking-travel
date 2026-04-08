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
import { FieldError } from "react-hook-form";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { SearchFieldHandle } from "./SearchForm";
import { Ref, useImperativeHandle, useRef, useState } from "react";

interface CalendarInputProps {
  placeholder: string;
  ref: Ref<SearchFieldHandle>;
  value?: Date;
  error?: FieldError;
  onSelect: (date: Date | undefined) => void;
}

export function CalendarInput({
  placeholder,
  value,
  onSelect,
  error,
  ref,
}: CalendarInputProps) {
  const { lang } = useParams();
  const locale = lang === "he" ? he : undefined;
  const displayValue = value ? value.toLocaleDateString(locale?.code) : "";
  const [open, setOpen] = useState(false);
  const triggerRef = useRef<HTMLInputElement>(null);

  useImperativeHandle(
    ref,
    () => ({
      focus() {
        setOpen(true);
        requestAnimationFrame(() => {
          triggerRef.current?.focus();
        });
      },
    }),
    [],
  );

  return (
    <div className="flex flex-col items-start">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Field>
            <InputGroup className="search-form-input px-0">
              <InputGroupInput
                aria-invalid={!!error}
                value={displayValue}
                placeholder={placeholder}
                className="text-start px-2"
                ref={triggerRef}
                readOnly
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
            locale={locale}
            numberOfMonths={2}
            selected={value}
            onSelect={(d) => {
              onSelect(d);
              setOpen(false);
            }}
            classNames={{
              today: "bg-brand/35",
              day_button:
                "text-navy data-[selected-single=true]:bg-brand data-[selected-single=true]:text-white",
            }}
          />
        </PopoverContent>
      </Popover>
      <ErrorDisplay>{error?.message}</ErrorDisplay>
    </div>
  );
}
