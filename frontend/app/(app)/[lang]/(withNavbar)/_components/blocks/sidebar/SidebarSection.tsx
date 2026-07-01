import { RichText } from "@payloadcms/richtext-lexical/react";
import type { SidebarSectionBlock } from "@/payload-types";
import { SidebarNav } from "./SidebarNav";
import { slugify } from "@/shared/utils/slugify";
import { SidebarToc } from "./SidebarToc";
import { clsx } from "clsx";

type SidebarSectionProps = {
  block: SidebarSectionBlock;
};

export function SidebarSection({ block }: SidebarSectionProps) {
  const sections = block.sections ?? [];

  if (sections.length === 0) return null;

  const type = block.type ?? "anchor";

  return (
    <div
      className={clsx("flex items-start gap-12 mb-20 lg:mx-auto", {
        "flex-col": type === "toc",
        "lg:w-2/3 mx-5": type === "anchor",
      })}
    >
      {type === "anchor" ? (
        <SidebarNav sections={sections} />
      ) : (
        <SidebarToc sections={sections} />
      )}

      <div className="flex flex-1 flex-col gap-12">
        {sections.map((section) => (
          <div
            key={section.id}
            id={slugify(section.title)}
            className="flex flex-col gap-2 scroll-mt-32"
          >
            <h4 className="py-3 type-h4 text-navy">{section.title}</h4>
            <div className="text-lg font-semibold leading-[1.7] tracking-tight text-gray-600 prose prose-headings:font-bold">
              <RichText data={section.content} />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
