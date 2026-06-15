import type { Form } from "@/payload-types";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import { ErrorDisplay } from "@/shared/components/ErrorDisplay";
import { ChevronDown } from "lucide-react";
import { Controller, FieldErrors, useForm } from "react-hook-form";
import { FormValues } from "./FormRenderer";

type FormField = NonNullable<Form["fields"]>[number];
const fieldClassName =
  "border-border-light bg-background h-[60px] rounded-xl px-6 type-paragraph text-text-secondary placeholder:text-text-secondary shadow-none";

function getFieldError(errors: FieldErrors<FormValues>, name: string) {
  const error = errors[name];

  return typeof error?.message === "string" ? error.message : undefined;
}

function getRequiredRule(required: boolean | null | undefined, label: string) {
  return required ? { value: true, message: label } : undefined;
}

export function FieldRenderer({
  control,
  errors,
  field,
  register,
  direction,
}: {
  control: ReturnType<typeof useForm<FormValues>>["control"];
  errors: FieldErrors<FormValues>;
  field: FormField;
  register: ReturnType<typeof useForm<FormValues>>["register"];
  direction: "ltr" | "rtl";
}) {
  const error = getFieldError(errors, field.name);
  const label = field.label ?? field.name;
  const required = getRequiredRule(field.required, label);
  const width = field.width ? `${field.width}%` : "100%";

  if (field.blockType === "textarea") {
    return (
      <div className="min-w-0" style={{ width }}>
        <Textarea
          aria-invalid={!!error}
          aria-label={label}
          className={cn(fieldClassName, "min-h-32.5 items-start py-6")}
          placeholder={label}
          {...register(field.name, { required })}
        />
        <ErrorDisplay>{error}</ErrorDisplay>
      </div>
    );
  }

  if (field.blockType === "select") {
    return (
      <div className="min-w-0" style={{ width }}>
        <Controller
          control={control}
          name={field.name}
          rules={{ required }}
          render={({ field: selectField }) => {
            const selectedLabel = field.options?.find(
              (option) => option.value === selectField.value,
            )?.label;

            return (
              <DropdownMenu dir={direction}>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    className={cn(
                      "flex w-full items-center justify-between border",
                      fieldClassName,
                      error && "border-destructive",
                    )}
                  >
                    <span>{selectedLabel ?? field.placeholder ?? label}</span>
                    <ChevronDown className="size-4 shrink-0 text-muted" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-(--radix-dropdown-menu-trigger-width)">
                  {field.options?.map((option) => (
                    <DropdownMenuItem
                      key={option.id ?? option.value}
                      onSelect={() => selectField.onChange(option.value)}
                    >
                      {option.label}
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            );
          }}
        />
        <ErrorDisplay>{error}</ErrorDisplay>
      </div>
    );
  }

  if (field.blockType === "checkbox") {
    return (
      <div className="min-w-0" style={{ width }}>
        <Controller
          control={control}
          name={field.name}
          rules={{ required }}
          render={({ field: checkboxField }) => (
            <Label className="flex cursor-pointer items-center gap-2 type-label text-foreground">
              <Checkbox
                aria-invalid={!!error}
                checked={Boolean(checkboxField.value)}
                className="size-4.5 border-[1.5px] border-border-light bg-white"
                onCheckedChange={(checked) =>
                  checkboxField.onChange(checked === true)
                }
              />
              <span>{label}</span>
            </Label>
          )}
        />
        <ErrorDisplay>{error}</ErrorDisplay>
      </div>
    );
  }

  const inputType = field.blockType === "email" ? "email" : "text";

  return (
    <div className="min-w-0" style={{ width }}>
      <Input
        aria-invalid={!!error}
        aria-label={label}
        className={fieldClassName}
        placeholder={label}
        type={inputType}
        {...register(field.name, {
          required,
          pattern:
            field.blockType === "email"
              ? {
                  value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                  message: label,
                }
              : undefined,
        })}
      />
      <ErrorDisplay>{error}</ErrorDisplay>
    </div>
  );
}
