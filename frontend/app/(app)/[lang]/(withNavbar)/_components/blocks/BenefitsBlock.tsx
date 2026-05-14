import Image from "next/image";
import type { BenefitsBlock } from "@/payload-types";
import type { Media } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SectionHeader } from "./SectionHeader";

interface Props {
  block: BenefitsBlock;
}

export function BenefitsBlock({ block }: Props) {
  if (!block.items || block.items.length === 0) return null;

  return (
    <section className="flex flex-col items-center gap-12 pb-30 w-2/3 mx-auto">
      <SectionHeader
        pillText={block.eyebrow}
        title={block.title ?? ""}
        subtitle={block.subtitle}
      />

      <div className="grid w-full grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {block.items.map((item) => {
          const media = item.image as Populated<typeof item.image>;
          return (
            <div
              key={item.id}
              className="flex flex-col items-center gap-6 rounded-2xl bg-white py-4 shadow-card"
            >
              {media && typeof media === "object" && (media as Media).url && (
                <Image
                  src={(media as Media).url!}
                  alt={(media as Media).alt ?? item.title}
                  width={63}
                  height={63}
                  className="h-16 w-auto object-contain"
                />
              )}
              <div className="flex flex-col items-center gap-2.5 text-center">
                <p className="type-h5 text-navy">{item.title}</p>
                {item.subtitle && (
                  <p className="type-h6 text-text-secondary font-normal">
                    {item.subtitle}
                  </p>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
