"use client";

import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { InputGroup, InputGroupAddon } from "@/components/ui/input-group";

interface Props {
  onSubmit: (phone: string) => void;
}

export function CustomerPhoneForm({ onSubmit }: Props) {
  const t = useTranslations("Login");
  const [phone, setPhone] = useState("");

  return (
    <div className="flex flex-col gap-3 w-full">
      <InputGroup
        dir="ltr"
        className="h-15 rounded-xl bg-background border-border-light"
      >
        <InputGroupAddon align="inline-start" className="gap-3 ps-6">
          <span className="type-paragraph text-text-secondary text-sm font-medium">
            +972
          </span>
          <div className="h-7.5 w-px bg-border-light/60" />
        </InputGroupAddon>
        <Input
          type="tel"
          placeholder={t("customer.phonePlaceholder")}
          value={phone}
          onChange={(e) => setPhone(e.target.value)}
          className="h-auto border-0 bg-transparent text-start type-paragraph text-text-secondary placeholder:text-text-secondary focus-visible:ring-0 pe-6"
        />
      </InputGroup>

      <Button
        variant="brand"
        className="w-full py-3.5 h-auto mt-3"
        onClick={() => onSubmit(phone)}
        disabled={!phone.trim()}
      >
        {t("customer.submit")}
      </Button>
    </div>
  );
}
