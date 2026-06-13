"use client";

import type { Form } from "@/payload-types";
import { Button } from "@/components/ui/button";
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
import { SuccessBadge } from "@/shared/components/UI/SuccessBadge";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { ChevronDown } from "lucide-react";
import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { Controller, FieldErrors, useForm } from "react-hook-form";
import { FieldRenderer } from "./FieldRenderer";
import { useDirection } from "@/shared/hooks/useDirection";

export type FormValues = Record<string, string | boolean>;

type PayloadFormRendererProps = {
  form: Form;
  title?: string | null;
  className?: string;
};

function buildDefaultValues(fields: Form["fields"]): FormValues {
  return (fields ?? []).reduce<FormValues>((values, field) => {
    values[field.name] =
      field.blockType === "checkbox"
        ? Boolean(field.defaultValue)
        : "defaultValue" in field
          ? (field.defaultValue ?? "")
          : "";

    return values;
  }, {});
}

export function PayloadFormRenderer({
  form,
  title,
  className,
}: PayloadFormRendererProps) {
  const dir = useDirection()
  
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [submitError, setSubmitError] = useState<string>();

  const {
    control,
    formState: { errors },
    handleSubmit,
    register,
  } = useForm<FormValues>({
    defaultValues: buildDefaultValues(form.fields),
  });

  const onSubmit = useCallback(
    async (data: FormValues) => {
      setIsLoading(true);
      setSubmitError(undefined);

      try {
        const response = await fetch("/api/form-submissions", {
          body: JSON.stringify({
            form: form.id,
            submissionData: Object.entries(data).map(([field, value]) => ({
              field,
              value: String(value ?? ""),
            })),
          }),
          headers: {
            "Content-Type": "application/json",
          },
          method: "POST",
        });

        if (!response.ok) {
          const result = await response.json().catch(() => null);
          setSubmitError(result?.errors?.[0]?.message);
          return;
        }

        setHasSubmitted(true);

        if (form.confirmationType === "redirect" && form.redirect?.url) {
          router.push(form.redirect.url);
        }
      } catch {
        setSubmitError(undefined);
      } finally {
        setIsLoading(false);
      }
    },
    [form.confirmationType, form.id, form.redirect?.url, router],
  );

  if (hasSubmitted && form.confirmationType === "message") {
    return (
      <SuccessBadge>
        {form.confirmationMessage ? (
          <RichText data={form.confirmationMessage} />
        ) : null}
      </SuccessBadge>
    );
  }

  return (
    <form
      className={cn(
        "flex w-full flex-col items-center gap-6",
        className,
      )}
      onSubmit={handleSubmit(onSubmit)}
    >
      {title && (
        <h4 className="type-h4 w-full text-start text-navy">{title}</h4>
      )}

      <div className="flex w-full flex-wrap gap-3">
        {form.fields?.map((field) => (
          <FieldRenderer
            key={field.id ?? field.name}
            control={control}
            errors={errors}
            field={field}
            register={register}
            direction={dir}
          />
        ))}
      </div>

      <div className="w-full">
        <Button
          type="submit"
          loading={isLoading}
          className="h-auto w-full rounded-[10px] bg-navy px-8 py-3.5 text-[15px] font-bold text-white shadow-subtle hover:bg-navy/90"
        >
          {form.submitButtonLabel}
        </Button>
        <ErrorDisplay>{submitError}</ErrorDisplay>
      </div>
    </form>
  );
}
