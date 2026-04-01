import { Page, SharedSectionRefBlock } from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SharedSectionRenderer } from "./sharedSections/sharedSectionRenderer";

export function BlocksRenderer({ blocks }: { blocks: Page["layout"] }) {
  return (
    <>
      {blocks?.map((block, index) => {
        switch (block.blockType) {
          case "sharedSectionRef":
            return (
              <SharedSectionRenderer
                key={index}
                section={
                  block.section as Populated<SharedSectionRefBlock["section"]>
                }
              />
            );
          default:
            return <div key={index}>[Unknown block: {block.blockType}]</div>;
        }
      })}
    </>
  );
}
