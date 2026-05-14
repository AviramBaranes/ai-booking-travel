import { useTranslations } from "next-intl";
import Image from "next/image";

interface SignalsDisplayProps {
  remainingCount: number;
  liveViewers: number;
}

export function SignalsDisplay({
  remainingCount,
  liveViewers,
}: SignalsDisplayProps) {
  const t = useTranslations("booking.results");

  return (
    <>
      <Image
        src="/assets/booking/signals/bell.gif"
        alt="Signal Bell"
        width={16}
        height={16}
      />
      <span className="type-label text-navy">
        {t("signals.remainingCount", {
          count: remainingCount,
        })}
      </span>
      <Image
        src="/assets/booking/signals/eye.gif"
        alt="Signal Eye"
        width={16}
        height={16}
      />
      <span className="type-label text-navy">
        {t("signals.liveViewers", {
          count: liveViewers,
        })}
      </span>
    </>
  );
}
