import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { BadgePercent, CalendarDays, User } from "lucide-react";
import { LogoutButton } from "./LogoutButton";
import { useTranslations } from "next-intl";
import useAuthStore from "@/shared/auth/authStore";
import { useParams, usePathname } from "next/navigation";
import { useState } from "react";
import { useDirection } from "@/shared/hooks/useDirection";

export function AuthenticatedDropdown() {
  const { lang } = useParams();
  const pathname = usePathname();
  const t = useTranslations("AuthDropdown");
  const user = useAuthStore((state) => state.user);
  const [open, setOpen] = useState(false);
  const dir = useDirection();

  if (!user || user.role === "admin") return null;

  const isAgent = user.role === "agent";

  const itemBase =
    "flex w-full items-center gap-2 px-3 py-3 text-[15px] font-medium transition-colors md:min-h-[71px] md:px-4 md:py-0 md:text-[16px]";

  function navItem(href: string) {
    const isActive = pathname === href || pathname.startsWith(href + "/");

    return isActive
      ? `${itemBase} bg-brand text-white border-b border-cars-border`
      : `${itemBase} text-navy border-b border-cars-border hover:bg-brand/30!`;
  }

  return (
    <DropdownMenu open={open} onOpenChange={setOpen} dir={dir}>
      <DropdownMenuTrigger asChild>
        <Button
          size="outline"
          variant="outline"
          className="px-2 py-1 lg:px-7 lg:py-4"
        >
          <User className="size-5 hidden lg:flex" />
          {t("greeting", { name: user.firstName })}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-45 lg:w-61.25 rounded-[12px] border border-border-light bg-white p-0 shadow-auth-dropdown overflow-hidden"
        align="center"
      >
        {/* Greeting header */}
        <div className="flex w-full items-center border-b border-cars-border px-3 py-4 text-[15px] font-bold text-navy md:min-h-18 md:px-4 md:py-0 md:text-[16px]">
          {t("greeting", { name: user.firstName })}
        </div>

        {/* Profile link */}
        {isAgent ? (
          <Link
            href={`/${lang}/price-offers`}
            className={navItem(`/${lang}/price-offers`)}
            onClick={() => setOpen(false)}
          >
            <BadgePercent
              className={`size-4 lg:size-6 shrink-0 ${pathname.startsWith(`/${lang}/price-offers`) ? "text-white" : "text-brand"}`}
            />
            <span>{t("priceOffers")}</span>
          </Link>
        ) : (
          <Link
            href={`/${lang}/profile`}
            className={navItem(`/${lang}/profile`)}
            onClick={() => setOpen(false)}
          >
            <User
              className={`size-4 lg:size-6 shrink-0 ${pathname.startsWith(`/${lang}/profile`) ? "text-white" : "text-brand"}`}
            />
            <span>{t("profile")}</span>
          </Link>
        )}

        {/* Reservations link */}
        <Link
          href={`/${lang}/reservations`}
          className={navItem(`/${lang}/reservations`)}
          onClick={() => setOpen(false)}
        >
          <CalendarDays
            className={`size-4 lg:size-6 shrink-0 ${pathname.startsWith(`/${lang}/reservations`) ? "text-white" : "text-brand"}`}
          />
          <span>{t("reservations")}</span>
        </Link>

        {/* Logout */}
        <LogoutButton
          buttonText={t("logout")}
          onLogout={() => setOpen(false)}
        />
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
