import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { CalendarDays, User } from "lucide-react";
import { LogoutButton } from "./LogoutButton";
import { useTranslations } from "next-intl";
import { useSession } from "next-auth/react";
import { useParams, usePathname } from "next/navigation";
import { useState } from "react";
import { useDirection } from "@/shared/hooks/useDirection";

export function AuthenticatedDropdown() {
  const { lang } = useParams();
  const pathname = usePathname();
  const t = useTranslations("AuthDropdown");
  const session = useSession();
  const [open, setOpen] = useState(false);
  const dir = useDirection();

  if (!session.data?.user || session.data.user.role === "admin") return null;

  const greetingKey =
    session.data.user.role === "agent" ? "helloAgent" : "helloCustomer";

  const itemBase =
    "flex items-center gap-2 px-4 min-h-[71px] w-full font-medium text-[16px] transition-colors";

  function navItem(href: string) {
    const isActive = pathname === href || pathname.startsWith(href + "/");
    return isActive
      ? `${itemBase} bg-brand text-white`
      : `${itemBase} text-navy border-b border-cars-border hover:bg-gray-50 hover:bg-brand/30!`;
  }

  return (
    <DropdownMenu open={open} onOpenChange={setOpen} dir={dir}>
      <DropdownMenuTrigger asChild>
        <Button size="outline" variant="outline">
          <User className="size-5" />
          {t(greetingKey)}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-61.25 rounded-[12px] border border-border-light bg-white p-0 shadow-auth-dropdown overflow-hidden"
        align="center"
      >
        {/* Greeting header */}
        <div className="flex items-center gap-2 px-4 min-h-18 w-full border-b border-cars-border font-bold text-[16px] text-navy">
          {t(greetingKey)}
        </div>

        {/* Profile link */}
        {/* <Link
          href={`/${lang}/profile`}
          className={navItem(`/${lang}/profile`)}
          onClick={() => setOpen(false)}
        >
          <User
            className={`size-6 shrink-0 ${pathname.startsWith(`/${lang}/profile`) ? "text-white" : "text-brand"}`}
          />
          <span>{t("profile")}</span>
        </Link> */}

        {/* Reservations link */}
        <Link
          href={`/${lang}/reservations`}
          className={navItem(`/${lang}/reservations`)}
          onClick={() => setOpen(false)}
        >
          <CalendarDays
            className={`size-6 shrink-0 ${pathname.startsWith(`/${lang}/reservations`) ? "text-white" : "text-brand"}`}
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
