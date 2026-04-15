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
import { DateRange } from "react-day-picker";

interface CalendarInputRangeProps {
  placeholder: string;
  ref: Ref<SearchFieldHandle>;
  value?: DateRange;
  error?: FieldError;
  onSelect: (date: DateRange | undefined) => void;
}

// CalenderInputRange is used only for a return date, showing only the return but showing a range while selecting
export function CalendarInputRange({
  placeholder,
  value,
  onSelect,
  error,
  ref,
}: CalendarInputRangeProps) {
  const { lang } = useParams();
  const locale = lang === "he" ? he : undefined;
  const displayValue = value?.to?.toLocaleDateString(locale?.code) ?? "";
  const [open, setOpen] = useState(false);
  const [hoverDate, setHoverDate] = useState<Date | undefined>(undefined);
  const triggerRef = useRef<HTMLInputElement>(null);

  useImperativeHandle(
    ref,
    () => ({
      focus() {
        setOpen(true);
      },
    }),
    [],
  );

  return (
    <div className="flex flex-col items-start">
      <Popover
        open={open}
        onOpenChange={(nextOpen) => {
          setOpen(nextOpen);
          if (!nextOpen) setHoverDate(undefined);
        }}
      >
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
            mode="range"
            showOutsideDays={false}
            locale={locale}
            numberOfMonths={2}
            selected={value}
            onSelect={(d) => {
              onSelect(d);
              setHoverDate(undefined);
              setOpen(false);
            }}
            classNames={{
              today: "bg-brand/35",

              day_button:
                "text-navy data-[selected-single=true]:bg-brand data-[selected-single=true]:text-white",
            }}
            disabled={(date) =>
              date < new Date() ||
              date > new Date(Date.now() + 365 * 24 * 60 * 60 * 1000) ||
              !value?.from ||
              date < value.from
            }
            previewFrom={value?.from && !value?.to ? value.from : undefined}
            previewTo={value?.from && !value?.to ? hoverDate : undefined}
            onPreviewDayEnter={(date) => setHoverDate(date)}
            onPreviewDayLeave={() => setHoverDate(undefined)}
          />
        </PopoverContent>
      </Popover>
      <ErrorDisplay>{error?.message}</ErrorDisplay>
    </div>
  );
}
