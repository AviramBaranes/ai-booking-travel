import { useCallback, useEffect, useRef } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

type UseDialogOpenFromQueryParams = {
  open: () => void;
};

export const OPEN_DIALOG_QUERY_KEY = "login";
export const OPEN_DIALOG_QUERY_VALUE = "open";

export function useDialogOpenFromQuery({ open }: UseDialogOpenFromQueryParams) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const hasOpenedFromQueryRef = useRef(false);
  const openRef = useRef(open);

  const matches =
    searchParams.get(OPEN_DIALOG_QUERY_KEY) === OPEN_DIALOG_QUERY_VALUE;

  useEffect(() => {
    openRef.current = open;
  }, [open]);

  useEffect(() => {
    if (matches) {
      if (!hasOpenedFromQueryRef.current) {
        hasOpenedFromQueryRef.current = true;
        openRef.current();
      }
      return;
    }

    hasOpenedFromQueryRef.current = false;
  }, [matches]);

  const clearQueryFlag = useCallback(() => {
    const hadFlag =
      searchParams.get(OPEN_DIALOG_QUERY_KEY) === OPEN_DIALOG_QUERY_VALUE;
    hasOpenedFromQueryRef.current = hadFlag;

    if (!hadFlag) {
      return;
    }

    const params = new URLSearchParams(searchParams.toString());
    params.delete(OPEN_DIALOG_QUERY_KEY);
    const next = params.toString();
    router.replace(next ? pathname + "?" + next : pathname, { scroll: false });
  }, [pathname, router, searchParams]);

  return { clearQueryFlag };
}
