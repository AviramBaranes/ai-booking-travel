import {
  BlogPost,
  Homepage,
  Page,
  SharedSectionRefBlock,
} from "@/payload-types";
import { Populated } from "@/shared/types/payload";
import { SharedSectionRenderer } from "./sharedSections/sharedSectionRenderer";
import { SidebarSection } from "./sidebar/SidebarSection";
import { SharedSectionWrapper } from "./sharedSections/SharedSectionWrapper";
import { RichText } from "@payloadcms/richtext-lexical/react";
import { FAQBlock } from "./faq/FaqBlock";
import { BenefitsBlock } from "./BenefitsBlock";
import { RelatedPosts } from "../posts/RelatedPosts";

export function BlocksRenderer({
  blocks,
  faqClassName,
  lang,
}: {
  blocks: Page["layout"] | Homepage["layout"];
  faqClassName?: string;
  lang?: string;
}) {
  return (
    <>
      {blocks?.map((block, index) => {
        switch (block.blockType) {
          case "sidebarSection":
            return <SidebarSection key={block.id} block={block} />;
          case "faq":
            return (
              <FAQBlock key={block.id} data={block} className={faqClassName} />
            );
          case "richText":
            return (
              <section
                className="lg:w-4/10 lg:mx-auto mx-5 prose prose-headings:font-bold max-w-none"
                key={block.id}
              >
                <RichText data={block.content} />
              </section>
            );
          case "sharedSectionRef":
            return (
              <SharedSectionWrapper key={block.id} overrides={block.overrides}>
                <SharedSectionRenderer
                  section={
                    block.section as Populated<SharedSectionRefBlock["section"]>
                  }
                />
              </SharedSectionWrapper>
            );
          case "benefits":
            return <BenefitsBlock key={block.id} block={block} />;
          case "related-posts":
            return (
              <section className="lg:w-2/3 mx-5 lg:mx-auto mb-10" key={block.id}>
                <RelatedPosts
                  showButton
                  pillText={block.eyebrow ?? ""}
                  posts={(block.relatedPosts as BlogPost[]) ?? []}
                  title={block.title ?? ""}
                  subtitle={block.subtitle ?? ""}
                  lang={lang ?? "he"}
                />
              </section>
            );
        }
      })}
    </>
  );
}
