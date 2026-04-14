import { booking } from "@/shared/client";
import { useDirection } from "@/shared/hooks/useDirection";
import clsx from "clsx";
import { useTranslations } from "next-intl";
import Image from "next/image";
import { SignalsDisplay } from "../../../_components/SignalsDisplay";

export function CarSignals({ vehicle }: { vehicle: booking.AvailableVehicle }) {
  const t = useTranslations("booking.results");
  const dir = useDirection();

  if (!vehicle.signals) return null;

  return (
    <div
      className={clsx("absolute top-0", {
        "right-0": dir === "ltr",
        "left-0": dir === "rtl",
      })}
    >
      {/* Gradient glow layer behind */}
      <div
        className={clsx(
          "absolute inset-0 bg-linear-to-r from-[rgba(53,112,181,0.3)] to-[rgba(236,25,138,0.3)] blur-lg",
          {
            "rounded-br-xl": dir === "rtl",
            "rounded-bl-xl": dir === "ltr",
          },
        )}
      />
      {/* Content layer on top */}
      <div
        className={clsx(
          "relative flex items-center gap-1.5 px-6 py-4 bg-[#fafafa] shadow-[0_4px_12px_0_rgba(53,112,181,0.14)]",
          {
            "rounded-tl-2xl rounded-br-xl": dir === "rtl",
            "rounded-tr-2xl rounded-bl-xl": dir === "ltr",
          },
        )}
      >
        <SignalsDisplay
          remainingCount={vehicle.signals.remainingCount}
          liveViewers={vehicle.signals.liveViewers}
        />
      </div>
    </div>
  );
}
