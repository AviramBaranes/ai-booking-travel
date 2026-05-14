"use client";

import { Button } from "@/components/ui/button";
import { ButtonGroup } from "@/components/ui/button-group";
import { cn } from "@/lib/utils";
import { MoreHorizontal } from "lucide-react";
import { useRouter } from "next/navigation";

const DEFAULT_ITEMS_PER_PAGE = 8;

const pageButtonClass =
  "h-9 w-9 rounded-none border border-[#e5e7eb] bg-white text-sm font-medium text-[#4a5565] hover:bg-gray-50 shadow-none";

export function PaginationButtons({
  total,
  page,
  searchParams,
  basePath,
  previousLabel,
  nextLabel,
  itemsPerPage = DEFAULT_ITEMS_PER_PAGE,
}: {
  total: number;
  page: number;
  searchParams: { toString: () => string };
  basePath: string;
  previousLabel: string;
  nextLabel: string;
  itemsPerPage?: number;
}) {
  const router = useRouter();
  const maxPage = Math.max(1, Math.ceil(total / itemsPerPage));

  function navigateTo(targetPage: number) {
    const nextQuery = new URLSearchParams(searchParams.toString());
    if (targetPage <= 1) {
      nextQuery.delete("page");
    } else {
      nextQuery.set("page", String(targetPage));
    }

    const queryString = nextQuery.toString();
    router.push(queryString ? `${basePath}?${queryString}` : basePath);
  }

  const windowStart = Math.max(1, page - 1);
  const windowEnd = Math.min(maxPage, page + 1);
  const windowPages: number[] = [];
  for (let i = windowStart; i <= windowEnd; i++) {
    windowPages.push(i);
  }
  const lastInWindow = windowPages[windowPages.length - 1];
  const showEllipsis = lastInWindow < maxPage - 1;
  const showLastPage = lastInWindow < maxPage;

  if (maxPage <= 1) return null;

  return (
    <div className="flex items-center justify-center mt-5" dir="ltr">
      <ButtonGroup
        aria-label="Pagination"
        className="overflow-hidden rounded-xl shadow-card"
      >
        <Button
          variant="outline"
          onClick={() => navigateTo(page - 1)}
          disabled={page <= 1}
          className={cn(pageButtonClass, "w-auto px-3 rounded-l-xl")}
        >
          {previousLabel}
        </Button>

        {windowPages.map((pageNumber) => (
          <Button
            key={pageNumber}
            variant="outline"
            onClick={() => navigateTo(pageNumber)}
            className={cn(
              pageButtonClass,
              pageNumber === page &&
                "bg-[#f9fafb] text-brand-blue hover:bg-[#f9fafb]",
            )}
          >
            {pageNumber}
          </Button>
        ))}

        {showEllipsis && (
          <Button
            variant="outline"
            disabled
            className={cn(
              pageButtonClass,
              "cursor-default disabled:opacity-100",
            )}
            tabIndex={-1}
          >
            <MoreHorizontal className="size-4" />
          </Button>
        )}

        {showLastPage && (
          <Button
            variant="outline"
            onClick={() => navigateTo(maxPage)}
            className={cn(
              pageButtonClass,
              maxPage === page &&
                "bg-[#f9fafb] text-brand-blue hover:bg-[#f9fafb]",
            )}
          >
            {maxPage}
          </Button>
        )}

        <Button
          variant="outline"
          onClick={() => navigateTo(page + 1)}
          disabled={page >= maxPage}
          className={cn(pageButtonClass, "w-auto px-3 rounded-r-xl")}
        >
          {nextLabel}
        </Button>
      </ButtonGroup>
    </div>
  );
}