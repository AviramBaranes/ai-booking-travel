import { booking } from "@/shared/client";
import { useTranslations } from "next-intl";
import Image from "next/image";

export function CarSignals({ vehicle }: { vehicle: booking.AvailableVehicle }) {
  const t = useTranslations("booking.results");

  if (!vehicle.signals) return null;

  return (
    <div className="absolute top-0 left-0">
      {/* Gradient glow layer behind */}
      <div className="absolute inset-0 rounded-br-xl bg-linear-to-r from-[rgba(53,112,181,0.3)] to-[rgba(236,25,138,0.3)] blur-lg" />
      {/* Content layer on top */}
      <div className="relative flex items-center gap-1.5 px-6 py-4 rounded-tl-2xl rounded-br-xl bg-[#fafafa] shadow-[0_4px_12px_0_rgba(53,112,181,0.14)]">
        <Image
          src="/assets/booking/signals/bell.gif"
          alt="Signal Bell"
          width={16}
          height={16}
        />
        <span className="type-label text-navy">
          {t("signals.remainingCount", {
            count: vehicle.signals.remainingCount,
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
            count: vehicle.signals.liveViewers,
          })}
        </span>
      </div>
    </div>
  );
}
