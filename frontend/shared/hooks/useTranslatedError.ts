import { useMemo } from "react";
import { isAppError } from "../api/AppError";
import { useTranslations } from "next-intl";

export function useTranslatedError(error: Error | null) {
  const tErrors = useTranslations("ApiErrors");

  return useMemo(() => {
    if (!error) return null;
    return tErrors(getErrorKey(error));
  }, [error]);
}

export function getErrorKey(error: Error): string {
  if (isAppError(error)) {
    return error.code;
  }

  return "internal_error";
}
