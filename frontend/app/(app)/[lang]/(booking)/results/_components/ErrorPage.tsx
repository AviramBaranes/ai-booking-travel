"use client";

import { useTranslations } from "next-intl";
import { SearchForm } from "../../../_components/home/SearchForm/SearchForm";

export default function ErrorResultPageContent() {
  const t = useTranslations("booking.results.error");

  return (
    <div className="mt-20">
      <div className="w-10/12 mx-auto">
        <SearchForm />
      </div>
      <div className="p-30 text-center">
        <h4 className="type-h4 text-navy">{t("noResults")}</h4>
      </div>
    </div>
  );
}
