import Image from "next/image";
import { TypedSection, Populated } from "@/shared/types/payload";
import type { Media } from "@/payload-types";
import { SectionHeader } from "../SectionHeader";

type StatsSectionProps = {
  section: TypedSection<"stats">;
};

export function StatsSection({ section }: StatsSectionProps) {
  const { pillText, title, subtitle, items } = section.stats;

  return (
    <section className="flex flex-col items-center gap-12 pb-20 w-8/12 mx-auto">
      <SectionHeader pillText={pillText} title={title} subtitle={subtitle} />

      {/* Stats cards */}
      {items && items.length > 0 && (
        <div className="flex w-full items-center justify-center gap-6">
          {items.map((item, index) => {
            const icon = item.icon as Populated<typeof item.icon>;

            return (
              <div key={item.id} className="contents">
                {/* Divider (between cards) */}
                {index > 0 && (
                  <div className="h-24 w-px shrink-0 bg-border-light" />
                )}

                {/* Card */}
                <div className="flex flex-1 flex-col items-center justify-center gap-1 rounded-xl py-6">
                  {/* Value + icon */}
                  <div className="relative">
                    <span className="relative z-10 text-3xl font-extrabold leading-tight tracking-tight text-foreground">
                      {item.value}
                    </span>
                    {icon?.url && (
                      <Image
                        src={icon.url}
                        alt=""
                        width={icon.width ?? 48}
                        height={icon.height ?? 48}
                        className="absolute -top-5 right-0 h-12 w-12 translate-x-1/2 object-contain"
                      />
                    )}
                  </div>

                  {/* Label */}
                  <span className="text-sm font-light text-foreground">
                    {item.label}
                  </span>

                  {/* Caption (optional) */}
                  {item.caption && (
                    <span className="text-xs font-light text-muted">
                      {item.caption}
                    </span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </section>
  );
}
