import { RichText } from "@payloadcms/richtext-lexical/react";
import type { SidebarSectionBlock } from "@/payload-types";
import { SidebarNav } from "./SidebarNav";
import { slugify } from "@/shared/lang/slugify";

type SidebarSectionProps = {
  block: SidebarSectionBlock;
};

export function SidebarSection({ block }: SidebarSectionProps) {
  const sections = block.sections ?? [];

  if (sections.length === 0) return null;

  return (
    <div className="flex items-start gap-12 w-2/3 mb-20 mx-auto">
      <SidebarNav sections={sections} />

      <div className="flex flex-1 flex-col gap-12 px-12">
        {sections.map((section) => (
          <div
            key={section.id}
            id={slugify(section.title)}
            className="flex flex-col gap-2 scroll-mt-32"
          >
            <h4 className="py-3 type-h4 text-navy">{section.title}</h4>
            <div className="text-lg font-semibold leading-[1.7] tracking-tight text-gray-600">
              <RichText data={section.content} />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
