import Link from "next/link";

type BlogPaginationProps = {
  lang: string;
  currentPage: number;
  totalPages: number;
};

export function BlogPagination({
  lang,
  currentPage,
  totalPages,
}: BlogPaginationProps) {
  if (totalPages <= 1) return null;

  const getHref = (page: number) => `/${lang}/blog/page/${page}`;

  const pages = getPaginationPages(currentPage, totalPages);

  return (
    <nav
      className="mt-12 flex items-center justify-center gap-2"
      aria-label="ניווט בין עמודי בלוג"
    >
      <PaginationLink
        href={getHref(currentPage - 1)}
        disabled={currentPage <= 1}
      >
        הקודם
      </PaginationLink>

      {pages.map((page, index) => {
        if (page === "...") {
          return (
            <span
              key={`dots-${index}`}
              className="flex h-10 min-w-10 items-center justify-center px-2 text-brand-blue/50"
            >
              ...
            </span>
          );
        }

        const isActive = page === currentPage;

        return (
          <Link
            key={page}
            href={getHref(page)}
            aria-current={isActive ? "page" : undefined}
            className={[
              "flex h-10 min-w-10 items-center justify-center rounded-lg border px-3 text-sm font-semibold transition",
              isActive
                ? "border-navy bg-navy text-white"
                : "border-border-light bg-white text-brand-blue hover:border-navy hover:text-navy",
            ].join(" ")}
          >
            {page}
          </Link>
        );
      })}

      <PaginationLink
        href={getHref(currentPage + 1)}
        disabled={currentPage >= totalPages}
      >
        הבא
      </PaginationLink>
    </nav>
  );
}

function PaginationLink({
  href,
  disabled,
  children,
}: {
  href: string;
  disabled: boolean;
  children: React.ReactNode;
}) {
  if (disabled) {
    return (
      <span className="flex h-10 items-center justify-center rounded-lg border border-border-light bg-gray-50 px-4 text-sm font-semibold text-brand-blue/35">
        {children}
      </span>
    );
  }

  return (
    <Link
      href={href}
      className="flex h-10 items-center justify-center rounded-lg border border-border-light bg-white px-4 text-sm font-semibold text-brand-blue transition hover:border-navy hover:text-navy"
    >
      {children}
    </Link>
  );
}

function getPaginationPages(
  currentPage: number,
  totalPages: number,
): Array<number | "..."> {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1);
  }

  if (currentPage <= 4) {
    return [1, 2, 3, 4, 5, "...", totalPages];
  }

  if (currentPage >= totalPages - 3) {
    return [
      1,
      "...",
      totalPages - 4,
      totalPages - 3,
      totalPages - 2,
      totalPages - 1,
      totalPages,
    ];
  }

  return [
    1,
    "...",
    currentPage - 1,
    currentPage,
    currentPage + 1,
    "...",
    totalPages,
  ];
}
