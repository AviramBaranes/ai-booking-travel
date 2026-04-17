import { clsx } from "clsx";

function Bone({ className }: { className?: string }) {
  return (
    <div className={clsx("animate-pulse rounded-md bg-white/50", className)} />
  );
}

export function SearchDataBannerDisplaySkeleton({
  dir = "rtl",
}: {
  dir?: "rtl" | "ltr";
}) {
  return (
    <section
      className="relative border-none w-full rounded-3xl bg-navy bg-cover bg-center bg-no-repeat shadow-card"
      style={{ backgroundImage: "url('/assets/booking/search-data-bg.png')" }}
    >
      <div className="flex items-center justify-between px-10 py-3">
        <div className="flex items-center gap-13">
          {/* Pickup */}
          <div className="flex flex-col gap-2 py-2">
            <Bone className="h-3 w-16" />
            <Bone className="h-5 w-40" />
            <Bone className="h-3 w-28" />
          </div>

          <div className="h-25 w-px bg-white/30" />

          {/* Dropoff */}
          <div className="flex flex-col gap-2 py-2">
            <Bone className="h-3 w-16" />
            <Bone className="h-5 w-40" />
            <Bone className="h-3 w-28" />
          </div>
        </div>
      </div>

      {/* Driver age */}
      <div
        className={clsx(
          "absolute bottom-0 flex items-center gap-2 border-t border-white px-3 py-4",
          {
            "left-0 rounded-bl-3xl rounded-tr-3xl border-r": dir === "rtl",
            "right-0 rounded-br-3xl rounded-tl-3xl border-l": dir === "ltr",
          },
        )}
      >
        <Bone className="h-4 w-4 rounded-full" />
        <Bone className="h-4 w-24" />
      </div>
    </section>
  );
}
