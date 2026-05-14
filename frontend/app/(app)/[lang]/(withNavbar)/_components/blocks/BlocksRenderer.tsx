import { Homepage, Page, SharedSectionRefBlock } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SharedSectionRenderer } from "./sharedSections/sharedSectionRenderer";
import { SidebarSection } from "./sidebar/SidebarSection";
import { SharedSectionWrapper } from "./sharedSections/SharedSectionWrapper";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { FAQBlock } from "./faq/FaqBlock";
import { BenefitsBlock } from "./BenefitsBlock";

export function BlocksRenderer({
  blocks,
  faqClassName,
}: {
  blocks: Page["layout"] | Homepage["layout"];
  faqClassName?: string;
}) {
  return (
    <>
      {blocks?.map((block, index) => {
        switch (block.blockType) {
          case "sidebarSection":
            return <SidebarSection key={index} block={block} />;
          case "faq":
            return (
              <FAQBlock key={index} data={block} className={faqClassName} />
            );
          case "richText":
            return <RichText key={index} data={block.content} />;
          case "sharedSectionRef":
            return (
              <SharedSectionWrapper key={index} overrides={block.overrides}>
                <SharedSectionRenderer
                  section={
                    block.section as Populated<SharedSectionRefBlock["section"]>
                  }
                />
              </SharedSectionWrapper>
            );
          case "benefits":
            return <BenefitsBlock key={index} block={block} />;
        }
      })}
    </>
  );
}
