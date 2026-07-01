import Image from "next/image";
import clsx from "clsx";
import { TypedSection, Populated } from "@/shared/types/payload";
import { SectionHeader } from "../SectionHeader";

type StatsSectionProps = {
  section: TypedSection<"stats">;
};

export function StatsSection({ section }: StatsSectionProps) {
  const { pillText, title, subtitle, items } = section.stats;

  return (
    <section className="mx-auto flex w-11/12 flex-col items-center gap-12 pb-20 lg:w-8/12">
      <SectionHeader pillText={pillText} title={title} subtitle={subtitle} />

      {items && items.length > 0 && (
        <div className="grid w-full grid-cols-2 lg:grid-cols-4">
          {items.map((item, index) => {
            const icon = item.icon as Populated<typeof item.icon>;

            return (
              <div
                key={item.id}
                className={clsx(
                  "relative flex flex-col items-center justify-center gap-1 rounded-xl py-6",
                  index % 2 === 1 &&
                    "before:absolute before:inset-y-6 before:inset-s-0 before:w-px before:bg-border-light",
                  index > 0 &&
                    "lg:before:absolute lg:before:inset-y-6 lg:before:inset-s-0 lg:before:w-px lg:before:bg-border-light",
                )}
              >
                <div className="relative">
                  <h3 className="relative z-10 type-h3 leading-tight tracking-tight text-foreground">
                    {item.value}
                  </h3>

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

                <h6 className="type-h6 text-foreground">{item.label}</h6>
              </div>
            );
          })}
        </div>
      )}
    </section>
  );
}
