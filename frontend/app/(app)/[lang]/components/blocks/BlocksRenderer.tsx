import { Page, SharedSectionRefBlock } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SharedSectionRenderer } from "./sharedSections/sharedSectionRenderer";
import { SidebarSection } from "./SidebarSection";
import { SharedSectionWrapper } from "./sharedSections/SharedSectionWrapper";

export function BlocksRenderer({ blocks }: { blocks: Page["layout"] }) {
  return (
    <>
      {blocks?.map((block, index) => {
        switch (block.blockType) {
          case "sidebarSection":
            return <SidebarSection key={index} block={block} />;
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
          default:
            return <div key={index}>[Unknown block: {block.blockType}]</div>;
        }
      })}
    </>
  );
}
