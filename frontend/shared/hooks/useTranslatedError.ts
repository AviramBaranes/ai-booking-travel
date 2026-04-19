import { useMemo } from "react";
import { isAppError } from "../api/AppError";
import { useTranslations } from "next-intl";

export function useTranslatedError(error: Error | null) {
  const tErrors = useTranslations("ApiErrors");

  return useMemo(() => {
    if (!error) return null;

    if (isAppError(error)) {
      return tErrors(error.code);
    }

    return tErrors("internal_error");
  }, [error]);
}
