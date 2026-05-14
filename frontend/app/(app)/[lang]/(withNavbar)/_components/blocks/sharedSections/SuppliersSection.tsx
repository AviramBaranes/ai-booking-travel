import Image from "next/image";
import { TypedSection, Populated } from "@/shared/types/payload";
import type { Media } from "@/payload-types";
import { SectionHeader } from "../SectionHeader";

type SuppliersSectionProps = {
  section: TypedSection<"suppliers">;
};

export function SuppliersSection({ section }: SuppliersSectionProps) {
  const { pillText, title, subtitle, logos } = section.suppliers;

  return (
    <section className="flex flex-col items-center gap-12 pb-20 w-10/12 mx-auto">
      <SectionHeader pillText={pillText} title={title} subtitle={subtitle} />

      {logos && logos.length > 0 && (
        <div className="grid w-full grid-cols-4 gap-px rounded-2xl bg-border-light overflow-hidden shadow-sm md:grid-cols-8">
          {logos.map((item) => {
            const media = item.logo as Populated<typeof item.logo>;
            return (
              <div
                key={item.id}
                className="flex items-center justify-center bg-white p-8"
              >
                {media?.url && (
                  <Image
                    src={media.url}
                    alt={media.alt ?? ""}
                    width={media.width ?? 120}
                    height={media.height ?? 60}
                    className="h-14 w-auto object-contain"
                  />
                )}
              </div>
            );
          })}
        </div>
      )}
    </section>
  );
}
