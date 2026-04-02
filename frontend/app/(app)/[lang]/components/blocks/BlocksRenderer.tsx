import { Page, SharedSectionRefBlock } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SharedSectionRenderer } from "./sharedSections/sharedSectionRenderer";
import { SidebarSection } from "./SidebarSection";
import { SharedSectionWrapper } from "./sharedSections/SharedSectionWrapper";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { FAQBlock } from "./FaqBlock";

export function BlocksRenderer({ blocks }: { blocks: Page["layout"] }) {
  return (
    <>
      {blocks?.map((block, index) => {
        switch (block.blockType) {
          case "sidebarSection":
            return <SidebarSection key={index} block={block} />;
          case "faq":
            return <FAQBlock key={index} data={block} />;
          case "richText":
            return <RichText key={index} data={block.content} />;
          case "sharedSectionRef":
            return (
              <SharedSectionWrapper overrides={block.overrides}>
                <SharedSectionRenderer
                  key={index}
                  section={
                    block.section as Populated<SharedSectionRefBlock["section"]>
                  }
                />
              </SharedSectionWrapper>
            );
        }
      })}
    </>
  );
}
