import "./SuppliersSection.css";
import Image from "next/image";
import { TypedSection, Populated } from "@/shared/types/payload";
import { SectionHeader } from "../SectionHeader";

type SuppliersSectionProps = {
  section: TypedSection<"suppliers">;
};

export function SuppliersSection({ section }: SuppliersSectionProps) {
  const { pillText, title, subtitle, logos } = section.suppliers;

  return (
    <section className="flex flex-col items-center gap-12 pb-20 w-full lg:w-10/12 lg:mx-auto">
      <SectionHeader pillText={pillText} title={title} subtitle={subtitle} />

      {logos && logos.length > 0 && (
        <div className="suppliers-logos-viewport">
          <div className="suppliers-logos-track">
            {logos.map((item) => {
              const media = item.logo as Populated<typeof item.logo>;

              if (!media?.url) return null;

              return (
                <div key={item.id} className="suppliers-logo-item">
                  <Image
                    src={media.url}
                    alt={media.alt ?? ""}
                    width={media.width ?? 88}
                    height={media.height ?? 88}
                    className="h-22 w-22 object-contain"
                  />
                </div>
              );
            })}
          </div>
        </div>
      )}
    </section>
  );
}
